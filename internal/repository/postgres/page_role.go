package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"gorm.io/gorm/clause"
)

func (r *PageRepository) CreatePageRole(ctx context.Context, createInput domain.PageRoleCreateInput) (*domain.PageRoleUser, *domain.Error) {
	tx, done := r.store.NewTransaction()
	defer done(nil)

	Email := strings.ToLower(createInput.Email)

	page := &model.Page{}
	if err := tx.DB().Where("pkid = ?", createInput.PagePkID).First(page).Error; err != nil {
		return nil, done(err)
	}

	var UserPkID *int64
	user := &model.User{}

	if err := tx.DB().Where("email = ?", Email).First(user).Error; err == nil {
		UserPkID = &user.Pkid
	}

	// New User page Role
	pageRole := model.PageRole{
		PagePkid: page.Pkid,
		Email:    Email,
		UserPkid: UserPkID,
		Role:     createInput.Role.String(),
	}

	if err := tx.DB().Create(&pageRole).Error; err != nil {
		return nil, done(err)
	}

	// Inherit Role to Child Pages
	childPath := pageutils.AppendPath(page.Path, strconv.FormatInt(page.Pkid, 10))
	childPages := []model.Page{}
	if err := tx.DB().Where("path LIKE ?", childPath+"%").Find(&childPages).Error; err != nil {
		return nil, done(err)
	}

	if len(childPages) != 0 {
		newRoles := sliceutils.Map(childPages, func(page model.Page) model.PageRole {
			return model.PageRole{
				PagePkid: page.Pkid,
				Email:    Email,
				UserPkid: UserPkID,
				Role:     domain.PageInherit.String(),
			}
		})

		// Ignore if child already has user role
		if err := tx.DB().Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&newRoles).Error; err != nil {
			return nil, done(err)
		}
	}

	return pageutils.TransformPageRoleModelToDomain(
		pageutils.PageRoleWithUser{
			PageRole: pageRole,
			User:     user,
		},
	), nil
}

func (r *PageRepository) GetPageRoleByEmail(
	ctx context.Context,
	pagePkID int64, email string,
) (*domain.PageRoleUser, *domain.Error) {
	var pageRole model.PageRole

	if err := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
		PagePkIDs: []int64{pagePkID},
		Emails:    []string{email},
	}).First(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	// If role inherit, get the role from parent page
	if pageRole.Role == domain.PageInherit.String() {
		var page model.Page
		if err := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; err != nil {
			return nil, domain.ErrDatabaseQuery
		}

		parentPkIDs := pageutils.PagePathToPkIDs(page.Path)

		var basePageRole model.PageRole
		query := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
			Emails:       []string{email},
			ExcludeRoles: []domain.PageRole{domain.PageInherit},
			PagePkIDs:    parentPkIDs,
			Preload: queryPageRolesPreloadOption{
				Page: true,
			},
		})
		if err := query.First(&basePageRole); err != nil {
			pageRole.Role = domain.PageViewer.String() // Default Role If Not Found inherit
		} else {
			pageRole.Role = basePageRole.Role
		}
	}

	return pageutils.TransformPageRoleModelToDomain(
		pageutils.PageRoleWithUser{
			PageRole: pageRole,
		},
	), nil
}

func (r *PageRepository) GetPageRoles(
	ctx context.Context,
	pagePkID int64,
) ([]domain.PageRoleUser, *domain.Error) {

	page := &model.Page{}
	if err := r.store.DB().Where("pkid = ?", pagePkID).First(page).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	pageRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
		PagePkIDs: []int64{pagePkID},
		Preload: queryPageRolesPreloadOption{
			User: true,
		},
	})
	if err != nil {
		return nil, err
	}

	inheritRoleEmails := sliceutils.Map(
		sliceutils.Filter(pageRoles, func(role PageRoleResult) bool {
			return role.Role == domain.PageInherit.String()
		}), func(role PageRoleResult) string {
			return role.Email
		})

	// Get Actual Page Role for Inherit Access
	parentPkIDs := pageutils.PagePathToPkIDs(page.Path)

	parentBasePageRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
		Emails:       inheritRoleEmails,
		ExcludeRoles: []domain.PageRole{domain.PageInherit},
		PagePkIDs:    parentPkIDs,
		Preload: queryPageRolesPreloadOption{
			Page: true,
		},
	})

	if err != nil {
		return nil, err
	}

	EmailRoleMap := make(map[string]PageRoleResult, len(parentBasePageRoles))
	for _, role := range parentBasePageRoles {
		EmailRoleMap[role.Email] = role
	}

	resultRoles := sliceutils.Map(pageRoles, func(role PageRoleResult) domain.PageRoleUser {
		transformRole := pageutils.PageRoleWithUser{
			PageRole: role.PageRole,
			Page:     role.Page,
			User:     role.User,
		}

		if role.Role == domain.PageInherit.String() {
			emailRole, ok := EmailRoleMap[role.Email]
			if !ok {
				transformRole.Role = domain.PageViewer.String() // Default Role If Not Found inherit
			} else {
				transformRole.Role = emailRole.Role
				inheritPage := emailRole.Page
				transformRole.InheritFromPage = pageutils.TransformPageModelToDomain(*inheritPage, nil, pageutils.PageBodyParams{}, nil)
			}
			role.Role = EmailRoleMap[role.Email].Role
		}
		return *pageutils.TransformPageRoleModelToDomain(transformRole)
	})

	return resultRoles, nil
}

func (r *PageRepository) UpdatePageRole(
	ctx context.Context,
	updateInput domain.PageRoleUpdateInput,
) *domain.Error {
	query := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
		PagePkIDs: []int64{updateInput.PagePkID},
		Emails:    []string{updateInput.Email},
	})
	if err := query.Model(&model.PageRole{}).Update("role", updateInput.Role.String()).Error; err != nil {
		return domain.ErrDatabaseMutation
	}

	return nil
}

func (r *PageRepository) DeletePageRole(
	ctx context.Context,
	updateInput domain.PageRoleDeleteInput,
) *domain.Error {
	tx, done := r.store.NewTransaction()
	defer done(nil)

	var page model.Page
	if err := tx.DB().Where("pkid = ?", updateInput.PagePkID).First(&page).Error; err != nil {
		return done(err)
	}

	if err := buildQueryPageRoles(tx.DB(), queryPageRolesParams{
		PagePkIDs: []int64{page.Pkid},
		Emails:    []string{updateInput.Email},
	}).Delete(&model.PageRole{}).Error; err != nil {
		return domain.ErrDatabaseMutation
	}

	// Remove All Inherit Roles from children
	childPath := pageutils.AppendPath(page.Path, strconv.FormatInt(page.Pkid, 10))
	childPages := []model.Page{}
	if err := tx.DB().Where("path LIKE ?", childPath+"%").Find(&childPages).Error; err != nil {
		return done(err)
	}
	if len(childPages) != 0 {
		q := buildQueryPageRoles(tx.DB(), queryPageRolesParams{
			PagePkIDs: sliceutils.Map(childPages, func(page model.Page) int64 {
				return page.Pkid
			}),
			Emails: []string{updateInput.Email},
			Roles:  []domain.PageRole{domain.PageInherit},
		})
		if err := q.Delete(&model.PageRole{}).Error; err != nil {
			return done(err)
		}
	}

	return nil
}

func (r *PageRepository) CheckPermission(ctx context.Context, input domain.PageRolePermissionCheckInput) (permissions domain.PageRolePermissions) {
	page := input.Page
	user := input.User

	if user == nil {
		if page.IsGeneralAccess {
			return getPermissionByRole(page.GeneralRole)
		}
		return permissions
	}

	// Author has all permissions
	if page.AuthorPkID != nil && user.PkID == *page.AuthorPkID {
		permissions.CanDelete = true
		permissions.CanEdit = true
		permissions.CanMove = true
		permissions.CanShare = true
		permissions.CanView = true

		return permissions
	}

	userRole, err := r.GetPageRoleByEmail(ctx, page.PkID, user.Email)

	// Direct Role
	if err == nil {
		permissions = getPermissionByRole(userRole.Role)
		return permissions
	}

	// Direct Role Not Found
	if page.IsGeneralAccess {
		return getPermissionByRole(page.GeneralRole)
	}

	return permissions
}

func getPermissionByRole(role domain.PageRole) (p domain.PageRolePermissions) {
	switch role {
	case domain.PageEditor:
		p.CanEdit = true
		p.CanView = true
		p.CanDownload = true
		p.CanShare = true
		p.CanDelete = true
		p.CanMove = true
	case domain.PageViewer:
		p.CanView = true
	case domain.PageInherit:
	}
	return p
}
