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

type DocRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewDocRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewDocRepository(params NewDocRepositoryParams) *DocRepository {
	return &DocRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *DocRepository) List(ctx context.Context, q domain.PageListQuery) ([]domain.Page, *domain.Error) {
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
		domainPages = append(domainPages, *pageutils.TransformPageModelToDomain(page, nil, nil))
	}
	return domainPages, nil
}

func (r *DocRepository) Update(ctx context.Context, pagePkID int64, updateInput domain.PageUpdateInput) (*domain.Page, *domain.Error) {
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
		nil,
	), nil
}

func (r *DocRepository) CreatePage(ctx context.Context, pageInput domain.PageInput) (*domain.Page, *domain.Error) {

	path := ""
	// get path
	if pageInput.ParentPagePkID != nil {
		var parentPage model.Page
		if err := r.store.DB().Where("pkid = ?", pageInput.ParentPagePkID).First(&parentPage).Error; err != nil {
			return nil, domain.NewErr("Parent Page not found", domain.BadRequestCode)
		}
		path = pageutils.AppendPath(parentPage.Path, parentPage.ID)
	}

	newPage := model.Page{
		Name:           pageInput.Name,
		CoverImage:     pageInput.CoverImage,
		OrgPkid:        &pageInput.OrganizationPkID,
		ParentPagePkid: pageInput.ParentPagePkID,
		ViewType:       pageInput.ViewType.String(),
		Path:           path,
	}
	if pageInput.Document.JsonContent == "" {
		pageInput.Document.JsonContent = "{}"
	}

	// Begin Tx
	tx, doneTx := r.store.NewTransaction()
	err := tx.DB().Create(&newPage).Error
	if err != nil {
		return nil, doneTx(err)
	}

	document := model.Document{
		JSONContent: &pageInput.Document.JsonContent,
		PagePkid:    newPage.Pkid,
	}

	rerr := tx.DB().Create(&document).Error
	if rerr != nil {
		return nil, doneTx(err)
	}

	doneTx(nil)
	// Commit Tx

	return pageutils.TransformPageModelToDomain(
		newPage,
		[]domain.Page{},
		pageutils.TransformDocModalToDomain(document),
	), nil
}

func (r *DocRepository) GetByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error) {
	var page model.Page
	if dbErr := r.store.DB().Where("id = ?", pageID).First(&page).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var childPages []model.Page
	if dbErr := r.store.DB().Where("parent_page_pkid = ?", page.Pkid).Find(&childPages).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var doc model.Document
	if dbErr := r.store.DB().Where("page_pkid = ?", page.Pkid).First(&doc).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	childPagesDomain := make([]domain.Page, len((childPages)))
	for i := 0; i < len(childPages); i++ {
		childPagesDomain[i] = *pageutils.TransformPageModelToDomain(childPages[i], nil, nil)
	}

	return pageutils.TransformPageModelToDomain(
		page,
		childPagesDomain,
		pageutils.TransformDocModalToDomain(doc),
	), nil
}

func (r *DocRepository) UpdateContent(ctx context.Context, pagePkID int64, content domain.DocumentInput) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	var doc model.Document
	if dbErr := r.store.DB().Where("page_pkid = ?", pagePkID).First(&doc).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}
	if content.JsonContent == "" {
		content.JsonContent = "{}"
	}
	doc.JSONContent = &content.JsonContent
	if dbErr := r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&doc).Error; dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return pageutils.TransformPageModelToDomain(
		page,
		nil,
		pageutils.TransformDocModalToDomain(doc),
	), nil
}

func (r *DocRepository) Archive(ctx context.Context, pagePkID int64) (*domain.Page, *domain.Error) {
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
		Where("path LIKE ?", descendantPath+"%").
		Select("ArchivedAt").
		Updates(model.Page{
			ArchivedAt: &now,
		}).Error; dbErr != nil {
		return nil, done(dbErr)
	}

	done(nil)

	return pageutils.TransformPageModelToDomain(page, nil, nil), nil
}

func (r *DocRepository) Move(ctx context.Context, pagePkID int64, parentPagePkID *int64) (*domain.Page, *domain.Error) {

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

	return pageutils.TransformPageModelToDomain(page, nil, nil), nil
}
