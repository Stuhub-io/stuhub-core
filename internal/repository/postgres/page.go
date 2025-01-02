package postgres

import (
	"context"
	"strconv"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"github.com/Stuhub-io/utils/userutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewPageRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewPageRepository(params NewPageRepositoryParams) *PageRepository {
	return &PageRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *PageRepository) List(
	ctx context.Context,
	q domain.PageListQuery,
	curUser *domain.User,
) ([]domain.Page, *domain.Error) {

	var results []PageResult

	query := buildPageQuery(preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Asset:     true,
		Doc:       true,
		Author:    true,
		PageRoles: true,
	}), q)

	if err := query.Find(&results).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var fallbackCurUserRole = domain.PageRestrict
	if q.ParentPagePkID != nil {
		parentRole, err := r.GetCurrentUserRole(ctx, *q.ParentPagePkID, curUser)
		if err != nil {
			return nil, err
		}
		fallbackCurUserRole = parentRole
	}

	domainPages := make([]domain.Page, 0, len(results))
	for _, result := range results {
		specificRole := fallbackCurUserRole

		if result.GeneralRole != domain.PageInherit.String() {
			specificRole = domain.PageRoleFromString(result.GeneralRole)
		} else {
			// if inherit, use parent's general role
			result.GeneralRole = specificRole.String()
		}
		if curUser != nil {
			userRole := sliceutils.Find(result.PageRoles, func(role model.PageRole) bool {
				return role.Email == curUser.Email
			})
			if userRole != nil {
				specificRole = domain.PageRoleFromString(userRole.Role)
			}
		}

		domainPage := pageutils.TransformPageModelToDomain(
			pageutils.PageModelToDomainParams{
				Page: &result.Page,
				PageBody: pageutils.PageBodyParams{
					Document: pageutils.TransformDocModelToDomain(result.Doc),
					Asset:    pageutils.TransformAssetModalToDomain(result.Asset),
				},
			},
		)

		permission := r.CheckPermission(ctx, domain.PageRolePermissionCheckInput{
			Page:     *domainPage,
			User:     curUser,
			PageRole: &specificRole,
		})

		domainPage.Permissions = &permission

		// Only return pages that user has permission to view
		if permission.CanView {
			domainPages = append(
				domainPages,
				*domainPage,
			)
		}
	}

	return domainPages, nil
}

func (r *PageRepository) Update(
	ctx context.Context,
	pagePkID int64,
	updateInput domain.PageUpdateInput,
) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr("Page not found", domain.BadRequestCode)
	}
	if updateInput.Name != nil && updateInput.Name != &page.Name {
		page.Name = *updateInput.Name
	}

	if updateInput.ViewType != nil && updateInput.ViewType.String() != page.ViewType {
		page.ViewType = updateInput.ViewType.String()
	}

	if updateInput.CoverImage != nil && *updateInput.CoverImage != page.CoverImage {
		page.CoverImage = *updateInput.CoverImage
	}

	dbErr := r.store.DB().
		Clauses(clause.Returning{}).
		Select("Name", "ViewType", "CoverImage").
		Save(&page).
		Error
	if dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page,
		},
	), nil
}

func (r *PageRepository) GetByID(
	ctx context.Context,
	pageID string,
	pagePkID *int64,
	detailOption domain.PageDetailOptions,
) (*domain.Page, *domain.Error) {

	var page PageResult

	query := preloadPageResult(r.store.DB().Model(&page), PreloadPageResultParams{
		Asset:  detailOption.Asset,
		Doc:    detailOption.Document,
		Author: detailOption.Author,
	})

	if pageID == "" && pagePkID == nil {
		return nil, domain.ErrBadParamInput
	}

	if pageID != "" {
		query = query.Where("id = ?", pageID)
	} else {
		query = query.Where("pkid = ?", *pagePkID)
	}

	if dbErr := query.First(&page).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	// Get Inherited Access
	var inheritPage *model.Page
	if page.GeneralRole == domain.PageInherit.String() {
		parentPkIDs := pageutils.PagePathToPkIDs(page.Path)
		query := buildPageQuery(r.store.DB(), domain.PageListQuery{
			OrgPkID:            *page.OrgPkid,
			ExcludeGeneralRole: []domain.PageRole{domain.PageInherit},
			PagePkIDs:          parentPkIDs,
		})
		if err := query.First(&inheritPage).Error; err != nil {
			page.GeneralRole = domain.PageRestrict.String()
		} else {
			page.GeneralRole = inheritPage.GeneralRole
		}
	}

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page.Page,
			PageBody: pageutils.PageBodyParams{
				Document: pageutils.TransformDocModelToDomain(page.Doc),
				Asset:    pageutils.TransformAssetModalToDomain(page.Asset),
				Author:   userutils.TransformUserModelToDomain(page.Author),
			},
			InheritFromPage: pageutils.TransformPageModelToDomain(
				pageutils.PageModelToDomainParams{
					Page: inheritPage,
				},
			),
		},
	), nil
}

func (r *PageRepository) Archive(
	ctx context.Context,
	pagePkID int64,
) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	now := time.Now()
	page.ArchivedAt = &now

	tx, done := r.store.NewTransaction()

	descendantPath := pageutils.AppendPath(page.Path, strconv.FormatInt(page.Pkid, 10))

	// Archive current page
	if dbErr := tx.DB().Clauses(clause.Locking{
		Strength: clause.LockingStrengthShare,
	}, clause.Returning{}).Select("ArchivedAt").Save(&page).Error; dbErr != nil {
		return nil, done(dbErr)
	}

	// Archive childrens
	if dbErr := tx.DB().Clauses(clause.Locking{
		Strength: clause.LockingStrengthShare,
	}, clause.Returning{}).
		Model(&model.Page{}).
		Where("path LIKE ? AND archived_at IS NULL", descendantPath+"%").
		Select("ArchivedAt").
		Updates(model.Page{
			ArchivedAt: &now,
		}).Error; dbErr != nil {
		return nil, done(dbErr)
	}

	done(nil)

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page,
		},
	), nil
}

func (r *PageRepository) Move(
	ctx context.Context,
	pagePkID int64,
	parentPagePkID *int64,
) (*domain.Page, *domain.Error) {

	var page model.Page

	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr("Page not found", domain.BadRequestCode)
	}

	oldPath := page.Path

	// Begin Tx
	tx, doneTx := r.store.NewTransaction()

	// get new path
	newPath := ""

	if parentPagePkID != nil {
		var parentPage model.Page
		if dbErr := tx.DB().Where("pkid = ?", parentPagePkID).First(&parentPage).Error; dbErr != nil {
			return nil, doneTx(dbErr)
		}
		newPath = pageutils.AppendPath(parentPage.Path, strconv.FormatInt(parentPage.Pkid, 10))
	}

	// update page path
	page.Path = newPath
	page.ParentPagePkid = parentPagePkID

	dbErr := tx.DB().Clauses(clause.Returning{}).Select("Path", "ParentPagePkid").Save(&page).Error

	if dbErr != nil {
		return nil, doneTx(dbErr)
	}

	descendantPath := pageutils.AppendPath(page.Path, strconv.FormatInt(page.Pkid, 10))
	descendantOldPath := pageutils.AppendPath(oldPath, page.ID)

	// batch update descendants
	bErr := tx.DB().
		Model(&model.Page{}).
		Where("path LIKE ?", descendantOldPath+"%").
		Update("path", gorm.Expr("replace(path, ?, ?)", descendantOldPath, descendantPath)).
		Error
	if bErr != nil {
		return nil, doneTx(bErr)
	}

	doneTx(nil)
	// Commit Tx

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page,
		},
	), nil
}

func (r *PageRepository) UpdateGeneralAccess(
	ctx context.Context,
	pagePkID int64,
	updateInput domain.PageGeneralAccessUpdateInput,
) (*domain.Page, *domain.Error) {

	page := model.Page{
		Pkid:        pagePkID,
		GeneralRole: updateInput.GeneralRole.String(),
	}

	if dbErr := r.store.DB().Clauses(clause.Returning{}).Select("GeneralRole").Save(&page).Error; dbErr != nil {
		return nil, domain.ErrUpdatePageGeneralAccess
	}

	return pageutils.TransformPageModelToDomain(
		pageutils.PageModelToDomainParams{
			Page: &page,
		},
	), nil
}

func (r *PageRepository) GetCurrentUserRole(ctx context.Context, pagePkID int64, curUser *domain.User) (domain.PageRole, *domain.Error) {
	var curUserRole = domain.PageRestrict
	// default to parent's general role
	parentPage, err := r.GetByID(ctx, "", &pagePkID, domain.PageDetailOptions{})
	if err != nil {
		return curUserRole, err
	}

	curUserRole = parentPage.GeneralRole
	// Check if user has direct role in parent page
	if curUser != nil {
		parentRoles, err := r.GetPageRoles(ctx, pagePkID)
		if err != nil {
			return curUserRole, err
		}
		curUserBaseRole := sliceutils.Find(parentRoles, func(role domain.PageRoleUser) bool {
			return role.Email == curUser.Email
		})

		if curUserBaseRole != nil {
			curUserRole = curUserBaseRole.Role
		}
	}
	return curUserRole, nil
}
