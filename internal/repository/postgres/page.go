package postgres

import (
	"context"
	"strconv"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/organizationutils"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"github.com/Stuhub-io/utils/userutils"
	"golang.org/x/exp/slices"
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

	// QUERY PAGES
	var results []PageResult
	query := buildPageQuery(preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Asset:     true,
		Doc:       true,
		Author:    true,
		PageRoles: true,
		PageStars: true,
	}), q)

	if err := query.Find(&results).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	// FILL UP PERMISSIONS
	userPageRoleLookupMap := make(map[int64]model.PageRole)
	generalPageRoleLookupMap := make(map[int64]domain.PageRole)

	// Find missing direct role for inherit direct roles
	missingPagePkIDs := make([]int64, 0, len(results))
	// Find missing page general role for inherit general roles
	missingGeneralPagePkIDs := make([]int64, 0, len(results))

	// Fill up PageRole Lookup Maps
	for _, result := range results {
		if result.GeneralRole != domain.PageInherit.String() {
			generalPageRoleLookupMap[result.Pkid] = domain.PageRoleFromString(result.GeneralRole)
		}

		if curUser != nil {
			curUserRole := sliceutils.Find(result.PageRoles, func(role model.PageRole) bool {
				return role.Email == curUser.Email
			})

			if curUserRole != nil && curUserRole.Role != domain.PageInherit.String() {
				userPageRoleLookupMap[result.Pkid] = *curUserRole
			}
		}
	}

	// Find missing Lookup Maps
	for _, result := range results {
		if result.GeneralRole == domain.PageInherit.String() {
			parentPkIDs := pageutils.PagePathToPkIDs(result.Path)
			for _, pkID := range parentPkIDs {
				if _, ok := generalPageRoleLookupMap[pkID]; !ok {
					missingGeneralPagePkIDs = append(missingGeneralPagePkIDs, pkID)
				}
			}
		}

		// Handle Direct Role
		if curUser != nil {
			// Find current user's direct page roles
			curUserRole := sliceutils.Find(result.PageRoles, func(role model.PageRole) bool {
				return role.Email == curUser.Email
			})

			if curUserRole != nil && curUserRole.Role == domain.PageInherit.String() {
				parentPkIDs := pageutils.PagePathToPkIDs(result.Path)
				for _, pkID := range parentPkIDs {
					if _, ok := userPageRoleLookupMap[pkID]; !ok {
						missingPagePkIDs = append(missingPagePkIDs, pkID)
					}
				}
			}

		}
	}

	// Get missing page general role lookup
	generalBasePages := []model.Page{}

	if buildPageQuery(r.store.DB(), domain.PageListQuery{
		OrgPkID:            q.OrgPkID,
		ExcludeGeneralRole: []domain.PageRole{domain.PageInherit},
		PagePkIDs:          missingGeneralPagePkIDs,
		IsAll:              true,
	}).Find(&generalBasePages).Error != nil {
		return nil, domain.ErrDatabaseQuery
	}

	for _, page := range generalBasePages {
		generalPageRoleLookupMap[page.Pkid] = domain.PageRoleFromString(page.GeneralRole)
	}

	// Fill missing direct roles
	if curUser != nil {
		// Find actual roles from missing direct role page pkids
		curUserIndirectRoles, err := queryPageRoles(r.store.DB(), queryPageRolesParams{
			Emails:       []string{curUser.Email},
			PagePkIDs:    missingPagePkIDs,
			ExcludeRoles: []domain.PageRole{domain.PageInherit},
		})

		if err != nil {
			return nil, err
		}

		for _, role := range curUserIndirectRoles {
			if _, ok := userPageRoleLookupMap[role.PagePkid]; !ok {
				userPageRoleLookupMap[role.PagePkid] = role.PageRole
			}
		}
	}

	findInheritGeneralRole := func(parentPkIds []int64) domain.PageRole {
		slices.Reverse(parentPkIds)
		for _, pkID := range parentPkIds {
			if role, ok := generalPageRoleLookupMap[pkID]; ok {
				return role
			}
		}
		return domain.PageRestrict
	}

	findCurPageRole := func(parentPkIDs []int64) *domain.PageRole {
		slices.Reverse(parentPkIDs)
		for _, pkID := range parentPkIDs {
			if role, ok := userPageRoleLookupMap[pkID]; ok {
				role := domain.PageRoleFromString(role.Role)
				return &role
			}
		}
		return nil
	}

	domainPages := make([]domain.Page, 0, len(results))

	for _, result := range results {
		// If inherit, find the base role
		if result.GeneralRole == domain.PageInherit.String() {
			result.GeneralRole = findInheritGeneralRole(pageutils.PagePathToPkIDs(result.Path)).String()
		}

		var directUserRole *domain.PageRole

		if curUser != nil {
			curPageRole := sliceutils.Find(result.PageRoles, func(role model.PageRole) bool {
				return role.Email == curUser.Email
			})

			// If direct role is inherit, find the base role
			if curPageRole != nil && curPageRole.Role == domain.PageInherit.String() {
				baseRole := findCurPageRole(pageutils.PagePathToPkIDs(result.Path))
				if baseRole != nil {
					directUserRole = baseRole
				}
			}

			// If direct role is not inherit, set to direct role
			if curPageRole != nil && curPageRole.Role != domain.PageInherit.String() {
				role := domain.PageRoleFromString(curPageRole.Role)
				directUserRole = &role
			}
		}

		var curUserPKID *int64 = nil
		if curUser != nil {
			curUserPKID = &curUser.PkID
		}

		domainPage := pageutils.TransformPageModelToDomain(
			pageutils.PageModelToDomainParams{
				Page: &result.Page,
				PageBody: pageutils.PageBodyParams{
					Document: pageutils.TransformDocModelToDomain(result.Doc),
					Asset:    pageutils.TransformAssetModalToDomain(result.Asset),
					Author:   userutils.TransformUserModelToDomain(result.Author),
				},
				PageStar: pageutils.GetPageStarByUserPkID(result.PageStars, curUserPKID),
			},
		)

		permission := r.CheckPermission(ctx, domain.PageRolePermissionCheckInput{
			Page:     *domainPage,
			User:     curUser,
			PageRole: directUserRole,
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
	actorPkID *int64,
) (*domain.Page, *domain.Error) {

	var page PageResult

	query := preloadPageResult(r.store.DB().Model(&page), PreloadPageResultParams{
		Asset:        detailOption.Asset,
		Doc:          detailOption.Document,
		Author:       detailOption.Author,
		Organization: detailOption.Organization,
		PageStars:    true,
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

	parentPagePkIDs := pageutils.PagePathToPkIDs(page.Path)

	// Get Actual General Role
	var inheritPage *model.Page
	if page.GeneralRole == domain.PageInherit.String() {
		query := buildPageQuery(r.store.DB(), domain.PageListQuery{
			OrgPkID:            page.OrgPkid,
			ExcludeGeneralRole: []domain.PageRole{domain.PageInherit},
			PagePkIDs:          parentPagePkIDs,
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
				Document:     pageutils.TransformDocModelToDomain(page.Doc),
				Asset:        pageutils.TransformAssetModalToDomain(page.Asset),
				Author:       userutils.TransformUserModelToDomain(page.Author),
				Organization: organizationutils.TransformOrganizationModelToDomain_Plain(page.Organization),
			},
			InheritFromPage: pageutils.TransformPageModelToDomain(
				pageutils.PageModelToDomainParams{
					Page: inheritPage,
				},
			),
			PageStar: pageutils.GetPageStarByUserPkID(page.PageStars, actorPkID),
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
	page.ParentPagePkid = nil

	tx, done := r.store.NewTransaction()

	descendantPath := pageutils.AppendPath(page.Path, strconv.FormatInt(page.Pkid, 10))

	// Archive current page, Move page to root
	if dbErr := tx.DB().Clauses(clause.Locking{
		Strength: clause.LockingStrengthShare, // FIXME: Need Locking ?
	}, clause.Returning{}).Select("ArchivedAt", "ParentPagePkid").Save(&page).Error; dbErr != nil {
		return nil, done(dbErr)
	}

	// Archive childrens
	if dbErr := tx.DB().Clauses(clause.Locking{
		Strength: clause.LockingStrengthShare, // FIXME: Need Locking ?
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
