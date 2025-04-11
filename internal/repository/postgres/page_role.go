package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"github.com/Stuhub-io/utils/userutils"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PageRepository) CreatePageRole(
	ctx context.Context,
	createInput domain.PageRoleCreateInput,
) (*domain.PageRoleUser, *domain.Error) {
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
	childPages := []PageResult{}

	if err := buildPageQuery(preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Author: true,
	}), domain.PageListQuery{
		PathBeginWith: childPath,
		IsAll:         true,
	}).Find(&childPages).Error; err != nil {
		return nil, done(err)
	}

	if len(childPages) != 0 {
		newRoles := sliceutils.Map(
			sliceutils.Filter(childPages, func(cp PageResult) bool { // Ignore if add email is already page's author
				return cp.Author.Email != Email
			}),
			func(page PageResult) model.PageRole {
				return model.PageRole{
					PagePkid: page.Pkid,
					Email:    Email,
					UserPkid: UserPkID,
					Role:     domain.PageInherit.String(),
				}
			},
		)

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
		return nil, domain.NewErr(err.Error(), domain.ErrInternalServerError.Code)
	}

	// If role inherit, get the role from parent page
	if pageRole.Role == domain.PageInherit.String() {
		var page model.Page
		if err := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; err != nil {
			return nil, domain.ErrDatabaseQuery
		}

		parentPkIDs := pageutils.PagePathToPkIDs(page.Path)
		slices.Reverse(parentPkIDs)

		var basePageRole model.PageRole
		query := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
			Emails:       []string{email},
			ExcludeRoles: []domain.PageRole{domain.PageInherit},
			PagePkIDs:    parentPkIDs,
			OrderBy:      buildOrderByValuesClause("page_pkid", parentPkIDs),
		})
		if err := query.First(&basePageRole).Error; err != nil {
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
	slices.Reverse(parentPkIDs)

	parentBasePageRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
		Emails:       inheritRoleEmails,
		ExcludeRoles: []domain.PageRole{domain.PageInherit},
		PagePkIDs:    parentPkIDs,
		Preload: queryPageRolesPreloadOption{
			Page: true,
		},
		OrderBy: buildOrderByValuesClause("page_pkid", parentPkIDs),
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
				transformRole.InheritFromPage = pageutils.TransformPageModelToDomain(
					pageutils.PageModelToDomainParams{
						Page: inheritPage,
					},
				)
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

func (r *PageRepository) GetPagesRole(
	ctx context.Context,
	input domain.PageRolePermissionBatchCheckInput,
) (permissions []domain.PageRolePermissionCheckInput, err *domain.Error) {
	user := input.User
	pages := input.Pages

	pageRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
		PagePkIDs: sliceutils.Map(pages, func(page domain.Page) int64 {
			return page.PkID
		}),
		Emails: []string{user.Email},
		Preload: queryPageRolesPreloadOption{
			Page: true,
		},
	})
	if err != nil {
		return nil, err
	}

	resultRoles := sliceutils.Map(
		pageRoles,
		func(role PageRoleResult) domain.PageRolePermissionCheckInput {
			pageRole := domain.PageRoleFromString(role.Role)
			return domain.PageRolePermissionCheckInput{
				User: user,
				Page: *pageutils.TransformPageModelToDomain(
					pageutils.PageModelToDomainParams{
						Page: role.Page,
					},
				),
				PageRole: &pageRole,
			}
		},
	)

	return resultRoles, nil
}

func (r *PageRepository) SyncPageRoleWithNewUser(
	ctx context.Context,
	user domain.User,
) *domain.Error {
	// Update User PkID in Page Roles
	if err := buildQueryPageRoles(r.store.DB(), queryPageRolesParams{
		Emails: []string{user.Email},
	}).Model(&model.PageRole{}).Update("user_pkid", user.PkID).Error; err != nil {
		return domain.ErrDatabaseMutation
	}
	return nil
}

func (r *PageRepository) CheckPermission(
	ctx context.Context,
	input domain.PageRolePermissionCheckInput,
) (permissions domain.PageRolePermissions) {
	page := input.Page
	user := input.User
	pageRoleUser := input.PageRole

	// General role
	if user == nil {
		return GetPermissionByRole(page.GeneralRole, false)
	}

	// User is Author
	if page.AuthorPkID != nil && user.PkID == *page.AuthorPkID {
		// Toggle all to True
		permissions.CanDelete = true
		permissions.CanEdit = true
		permissions.CanMove = true
		permissions.CanShare = true
		permissions.CanView = true
		permissions.CanDownload = true

		return permissions
	}

	if pageRoleUser != nil {
		permissions = GetPermissionByRole(*pageRoleUser, true)
		return permissions
	}

	// Direct Role Not Found
	return GetPermissionByRole(page.GeneralRole, true)
}

// Exclude Inherit role instead.
func GetPermissionByRole(role domain.PageRole, isAuthenticated bool) (p domain.PageRolePermissions) {
	if role != domain.PageRestrict && !isAuthenticated {
		p.CanView = true
		p.CanDownload = true
		return p
	}
	switch role {
	case domain.PageRestrict:
		return p
	case domain.PageEditor:
		p.CanEdit = true
		p.CanView = true
		p.CanDownload = true
		p.CanShare = true
		p.CanDelete = true
		p.CanMove = true
	case domain.PageViewer:
		p.CanView = true
		p.CanDownload = true
	case domain.PageInherit:
	}
	return p
}

func (r *PageRepository) CreatePageAccessRequest(ctx context.Context, input domain.PageRoleRequestCreateInput) (*domain.PageRoleRequestLog, *domain.Error) {

	var UserPkID *int64
	user := &model.User{}

	if err := r.store.DB().Where("email = ?", input.Email).First(user).Error; err == nil {
		UserPkID = &user.Pkid
	}

	pageRoleRequest := model.PagePermissionRequestLog{
		PagePkid: input.PagePkID,
		Email:    input.Email,
		Status:   domain.PRSLPending.String(),
		UserPkid: UserPkID,
	}

	if err := r.store.DB().Create(&pageRoleRequest).Error; err != nil {
		return nil, domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}

	return pageutils.TransformPagePermissionRequestLogToDomain(pageutils.PagePermissionRequestLogToDomainParams{
		Model: &pageRoleRequest,
	}), nil
}

type PageRoleRequestLogResults struct {
	model.PagePermissionRequestLog
	User *model.User `gorm:"foreignKey:user_pkid"`
}

func (r *PageRepository) ListPageAccessRequestByPagePkID(ctx context.Context, q domain.PageRoleRequestLogQuery) ([]domain.PageRoleRequestLog, *domain.Error) {
	// Write build query + preload utils for this
	requests := []PageRoleRequestLogResults{}

	query := buildPageAccessRequestQuery(r.store.DB().Preload("User"), q)

	if err := query.Order("created_at desc").Find(&requests).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	// Remove Duplicate Requests, get latest request

	listedEmail := make(map[string]bool, len(requests))
	p := sliceutils.Map(
		sliceutils.Filter(requests, func(r PageRoleRequestLogResults) bool {
			if _, ok := listedEmail[r.Email]; ok {
				return false
			}
			listedEmail[r.Email] = true
			return true
		}),
		func(r PageRoleRequestLogResults) domain.PageRoleRequestLog {
			return *pageutils.TransformPagePermissionRequestLogToDomain(pageutils.PagePermissionRequestLogToDomainParams{
				Model: &r.PagePermissionRequestLog,
				User:  userutils.TransformUserModelToDomain(r.User),
			})
		})
	return p, nil
}

func (r *PageRepository) UpdatePageAccessRequestStatus(ctx context.Context, q domain.PageRoleRequestLogQuery, status domain.PageRoleRequestLogStatus) *domain.Error {
	query := buildPageAccessRequestQuery(r.store.DB().Model(&PageRoleRequestLogResults{}), q)
	if err := query.Update("status", status.String()).Error; err != nil {
		return domain.ErrDatabaseMutation
	}
	return nil
}

func buildPageAccessRequestQuery(tx *gorm.DB, q domain.PageRoleRequestLogQuery) *gorm.DB {
	query := tx

	if len(q.PagePkIDs) != 0 {
		if len(q.PagePkIDs) == 1 {
			query = query.Where("page_pkid = ?", q.PagePkIDs[0])
		} else {
			query = query.Where("page_pkid IN ?", q.PagePkIDs)
		}
	}

	if len(q.Status) != 0 {
		query = query.Where("status IN ?", q.Status)
	}

	if len(q.Emails) != 0 {
		if len(q.Emails) == 1 {
			query = query.Where("email = ?", q.Emails)
		} else {
			query = query.Where("email IN ?", q.Emails)
		}
	}

	return query
}
