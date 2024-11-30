package postgres

import (
	"context"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
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

func (r *PageRepository) List(ctx context.Context, q domain.PageListQuery) ([]domain.Page, *domain.Error) {
	var pages []model.Page
	query := r.store.DB().Where("org_pkid = ?", q.OrgPkID)

	// Filter by Archived
	if q.IsArchived != nil {
		query = query.Where("archived_at IS NULL = ?", !*q.IsArchived)
	}

	// Filter By Parent Page
	if q.ParentPagePkID != nil {
		query = query.Where("parent_page_pkid = ?", *q.ParentPagePkID)
	}

	if len(q.ViewTypes) > 0 {
		query = query.Where("view_type IN ?", q.ViewTypes)
	}

	query = query.Order("created_at desc").Offset(q.Offset).Limit(q.Limit)

	err := query.Find(&pages).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	domainPages := make([]domain.Page, 0, len(pages))
	for _, page := range pages {
		domainPages = append(domainPages, *pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}))
	}
	return domainPages, nil
}

func (r *PageRepository) Update(ctx context.Context, pagePkID int64, updateInput domain.PageUpdateInput) (*domain.Page, *domain.Error) {
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

	dbErr := r.store.DB().Clauses(clause.Returning{}).Select("Name", "ViewType", "CorverImage").Save(&page).Error
	if dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return pageutils.TransformPageModelToDomain(
		page,
		nil,
		pageutils.PageBodyParams{},
	), nil
}

func (r *PageRepository) GetByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error) {
	var page model.Page
	if dbErr := r.store.DB().Where("id = ?", pageID).First(&page).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var childPages []model.Page
	if dbErr := r.store.DB().Where("parent_page_pkid = ?", page.Pkid).Find(&childPages).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var pageInstance struct {
		Doc   *model.Document
		Asset *model.Asset
	}
	switch domain.PageViewFromString(page.ViewType) {
	case domain.PageViewTypeDoc:
		var doc model.Document
		if dbErr := r.store.DB().Where("page_pkid = ?", page.Pkid).First(&doc).Error; dbErr != nil {
			return nil, domain.ErrDatabaseQuery
		}
		pageInstance.Doc = &doc
	case domain.PageViewTypeAsset:
		var asset model.Asset
		if dbErr := r.store.DB().Where("page_pkid = ?", page.Pkid).First(&asset).Error; dbErr != nil {
			return nil, domain.ErrDatabaseQuery
		}
		pageInstance.Asset = &asset
	case domain.PageViewTypeFolder:
		// Do nothing
	}

	childPagesDomain := make([]domain.Page, len((childPages)))
	for i := 0; i < len(childPages); i++ {
		childPagesDomain[i] = *pageutils.TransformPageModelToDomain(childPages[i], nil, pageutils.PageBodyParams{})
	}

	return pageutils.TransformPageModelToDomain(
		page,
		childPagesDomain,
		pageutils.PageBodyParams{
			Document: pageutils.TransformDocModelToDomain(pageInstance.Doc),
			Asset:    pageutils.TransformAssetModalToDomain(pageInstance.Asset),
		},
	), nil
}

func (r *PageRepository) Archive(ctx context.Context, pagePkID int64) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	now := time.Now()
	page.ArchivedAt = &now

	tx, done := r.store.NewTransaction()

	descendantPath := pageutils.AppendPath(page.Path, page.ID)

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

	return pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}), nil
}

func (r *PageRepository) Move(ctx context.Context, pagePkID int64, parentPagePkID *int64) (*domain.Page, *domain.Error) {

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
		newPath = pageutils.AppendPath(parentPage.Path, parentPage.ID)
	}

	// update page path
	page.Path = newPath
	page.ParentPagePkid = parentPagePkID

	dbErr := tx.DB().Clauses(clause.Returning{}).Select("Path", "ParentPagePkid").Save(&page).Error

	if dbErr != nil {
		return nil, doneTx(dbErr)
	}

	descendantPath := pageutils.AppendPath(page.Path, page.ID)
	descendantOldPath := pageutils.AppendPath(oldPath, page.ID)

	// batch update descendants
	bErr := tx.DB().Model(&model.Page{}).Where("path LIKE ?", descendantOldPath+"%").Update("path", gorm.Expr("replace(path, ?, ?)", descendantOldPath, descendantPath)).Error
	if bErr != nil {
		return nil, doneTx(bErr)
	}

	doneTx(nil)
	// Commit Tx

	return pageutils.TransformPageModelToDomain(page, nil, pageutils.PageBodyParams{}), nil
}
