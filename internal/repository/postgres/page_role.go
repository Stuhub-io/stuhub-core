package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

func (r *PageRepository) CreatePageRole(ctx context.Context, createInput domain.PageRoleCreateInput) (*domain.PageRoleUser, *domain.Error) {

	Email := strings.ToLower(createInput.Email)
	var UserPkID *int64
	user := &model.User{}

	if err := r.store.DB().Where("email = ?", Email).First(user).Error; err == nil {
		UserPkID = &user.Pkid
	}

	pageRole := model.PageRole{
		PagePkid: createInput.PagePkID,
		Email:    Email,
		UserPkid: UserPkID,
		Role:     createInput.Role.String(),
	}

	if err := r.store.DB().Create(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	// FIXME: Create an INHERIT type Role for all descendant Pages for new email -- if not exist

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
		PkIDs:  []int64{pagePkID},
		Emails: []string{email},
	}).First(&pageRole).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	// FIXME: Get Actual Page Role From Parent for Inherit Access

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
		PkIDs: []int64{pagePkID},
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
	parentPkIDs := sliceutils.Map(strings.Split(page.Path, "/"), func(pkid string) int64 {
		parsedPkID, err := strconv.ParseInt(pkid, 10, 64)
		if err != nil {
			return -1
		}
		return parsedPkID
	})

	parentBasePageRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
		Emails:       inheritRoleEmails,
		ExcludeRoles: []domain.PageRole{domain.PageInherit},
		PkIDs:        parentPkIDs,
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
		PkIDs:  []int64{updateInput.PagePkID},
		Emails: []string{updateInput.Email},
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

	if err := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
		PkIDs:  []int64{updateInput.PagePkID},
		Emails: []string{updateInput.Email},
	}).Delete(&model.PageRole{}).Error; err != nil {
		return domain.ErrDatabaseMutation
	}

	// FIXME: DELETE ALL INHERIT ROLES FOR DESCENDANT PAGES

	return nil
}
