package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
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
) ([]domain.Page, *domain.Error) {
	var results []PageResult
	query := r.preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Asset:  true,
		Doc:    true,
		Author: true,
	}).
		Where("org_pkid = ?", q.OrgPkID)

	if q.IsArchived != nil {
		if *q.IsArchived {
			query = query.Where("pages.archived_at IS NOT NULL")
		} else {
			query = query.Where("pages.archived_at IS NULL")
		}
	}
	if !q.IsAll {
		if q.ParentPagePkID != nil {
			query = query.Where("pages.parent_page_pkid = ?", *q.ParentPagePkID)
		} else {
			query = query.Where("pages.parent_page_pkid IS NULL")
		}
	}
	if (len(q.ViewTypes)) > 0 {
		query = query.Where("pages.view_type IN ?", q.ViewTypes)
	}

	query = query.Order("pages.updated_at desc").Offset(q.Offset).Limit(q.Limit)

	if err := query.Find(&results).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	domainPages := make([]domain.Page, 0, len(results))
	for _, result := range results {
		domainPages = append(
			domainPages,
			*pageutils.TransformPageModelToDomain(result.Page, nil, pageutils.PageBodyParams{
				Document: pageutils.TransformDocModelToDomain(result.Doc),
				Asset:    pageutils.TransformAssetModalToDomain(result.Asset),
			}, nil),
		)
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
		page,
		nil,
		pageutils.PageBodyParams{},
		nil,
	), nil
}

func (r *PageRepository) GetByID(
	ctx context.Context,
	pageID string,
	pagePkID *int64,
) (*domain.Page, *domain.Error) {
	var page PageResult

	query := r.preloadPageResult(r.store.DB().Model(&page), PreloadPageResultParams{
		Asset:  true,
		Doc:    true,
		Author: true,
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

	var childPages []PageResult
	if dbErr := r.preloadPageResult(r.store.DB(), PreloadPageResultParams{
		Asset: true,
		Doc:   true,
	}).Where("parent_page_pkid = ?", page.Pkid).Order("pages.updated_at desc").Find(&childPages).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	childPagesDomain := make([]domain.Page, len((childPages)))
	for i := 0; i < len(childPages); i++ {
		childPagesDomain[i] = *pageutils.TransformPageModelToDomain(childPages[i].Page, nil, pageutils.PageBodyParams{
			Document: pageutils.TransformDocModelToDomain(childPages[i].Doc),
			Asset:    pageutils.TransformAssetModalToDomain(childPages[i].Asset),
		}, nil)
	}

	return pageutils.TransformPageModelToDomain(
		page.Page,
		childPagesDomain,
		pageutils.PageBodyParams{
			Document: pageutils.TransformDocModelToDomain(page.Doc),
			Asset:    pageutils.TransformAssetModalToDomain(page.Asset),
			Author:   userutils.TransformUserModelToDomain(page.Author),
		},
		nil,
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

	return pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}, nil), nil
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

	return pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}, nil), nil
}

func (r *PageRepository) UpdateGeneralAccess(
	ctx context.Context,
	pagePkID int64,
	updateInput domain.PageGeneralAccessUpdateInput,
) (*domain.Page, *domain.Error) {
	page := model.Page{
		Pkid:            pagePkID,
		IsGeneralAccess: updateInput.IsGeneralAccess,
		GeneralRole:     updateInput.GeneralRole.String(),
	}

	if dbErr := r.store.DB().Clauses(clause.Returning{}).Select("IsGeneralAccess", "GeneralRole").Save(&page).Error; dbErr != nil {
		fmt.Print(dbErr)
		return nil, domain.ErrUpdatePageGeneralAccess
	}

	return pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}, nil), nil
}
