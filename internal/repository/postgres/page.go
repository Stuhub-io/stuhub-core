package postgres

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/docutils"
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

	var domainPages []domain.Page = make([]domain.Page, 0, len(pages))
	for _, page := range pages {
		domainPages = append(domainPages, *docutils.TransformPageModelToDomain(page, nil, nil))
	}
	return domainPages, nil
}

func (r *DocRepository) Update(ctx context.Context, pagePkID int64, updateInput domain.PageUpdateInput) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr("Page not found", domain.BadRequestCode)
	}
	if updateInput.Name != nil {
		page.Name = *updateInput.Name
	}
	if updateInput.ViewType != nil {
		page.ViewType = updateInput.ViewType.String()
	}
	if updateInput.ParentPagePkID != nil {
		page.ParentPagePkid = updateInput.ParentPagePkID
	}
	if updateInput.CoverImage != nil {
		page.CoverImage = *updateInput.CoverImage
	}

	dbErr := r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&page).Error
	if dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return docutils.TransformPageModelToDomain(
		page,
		nil,
		nil,
	), nil
}

func (r *DocRepository) CreatePage(ctx context.Context, pageInput domain.PageInput) (*domain.Page, *domain.Error) {
	newPage := model.Page{
		Name:           pageInput.Name,
		CoverImage:     pageInput.CoverImage,
		OrgPkid:        &pageInput.OrganizationPkID,
		ParentPagePkid: pageInput.ParentPagePkID,
		ViewType:       pageInput.ViewType.String(),
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

	return docutils.TransformPageModelToDomain(
		newPage,
		[]domain.Page{},
		docutils.TransformDocModalToDomain(document),
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
		childPagesDomain[i] = *docutils.TransformPageModelToDomain(childPages[i], nil, nil)
	}

	return docutils.TransformPageModelToDomain(
		page,
		childPagesDomain,
		docutils.TransformDocModalToDomain(doc),
	), nil
}

func (r *DocRepository) UpdateContent(ctx context.Context, pagePkID int64, jsonContent string) (*domain.Page, *domain.Error) {
	var page = model.Page{}
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	var doc model.Document
	if dbErr := r.store.DB().Where("page_pkid = ?", pagePkID).First(&doc).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}
	if jsonContent == "" {
		jsonContent = "{}"
	}
	doc.JSONContent = &jsonContent
	if dbErr := r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&doc).Error; dbErr != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return docutils.TransformPageModelToDomain(
		page,
		nil,
		docutils.TransformDocModalToDomain(doc),
	), nil
}

func (r *DocRepository) Archive(ctx context.Context, pagePkID int64) (*domain.Page, *domain.Error) {
	return nil, nil
}
