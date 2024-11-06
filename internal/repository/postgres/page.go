package postgres

// import (
// 	"context"

// 	"time"

// 	"github.com/Stuhub-io/config"
// 	"github.com/Stuhub-io/core/domain"
// 	"github.com/Stuhub-io/core/ports"
// 	store "github.com/Stuhub-io/internal/repository"
// 	"github.com/Stuhub-io/internal/repository/model"
// 	"github.com/Stuhub-io/utils/pageutils"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm/clause"
// )

// type PageRepository struct {
// 	store *store.DBStore
// 	cfg   config.Config
// }

// type NewPageRepositoryParams struct {
// 	Cfg   config.Config
// 	Store *store.DBStore
// }

// func NewPageRepository(params NewPageRepositoryParams) ports.PageRepository {
// 	return &PageRepository{
// 		store: params.Store,
// 		cfg:   params.Cfg,
// 	}
// }

// func (r *PageRepository) GetPagesBySpacePkID(ctx context.Context, spacePkID int64, excludeArchived bool) ([]domain.Page, *domain.Error) {
// 	var pages []model.Page
// 	query := r.store.DB().Where("space_pkid = ?", spacePkID)
// 	if excludeArchived {
// 		query = query.Where("archived_at IS NULL")
// 	}
// 	err := query.Order("created_at desc").Find(&pages).Error
// 	if err != nil {
// 		return nil, domain.ErrDatabaseQuery
// 	}

// 	var domainPages []domain.Page
// 	for _, page := range pages {
// 		domainPages = append(domainPages, *pageutils.TransformPageModelToDomain(page, nil))
// 	}

// 	return domainPages, nil
// }

// func (r *PageRepository) DeletePageByPkID(ctx context.Context, pagePkID int64, userPkID int64) (*domain.Page, *domain.Error) {
// 	var page model.Page
// 	isDeleted := r.store.DB().Where("pkid = ?", pagePkID).Delete(&page).Error
// 	if isDeleted != nil {
// 		return nil, domain.ErrDatabaseDelete
// 	}

// 	return pageutils.TransformPageModelToDomain(page, nil), nil
// }

// func (r *PageRepository) GetPageByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error) {
// 	var page model.Page
// 	err := r.store.DB().Where("id = ?", pageID).First(&page).Error
// 	if err != nil {
// 		return nil, domain.ErrDatabaseQuery
// 	}
// 	var childPages []model.Page
// 	rerr := r.store.DB().Where("parent_page_pkid = ?", page.Pkid).Find(&childPages).Error
// 	if rerr != nil {
// 		return nil, domain.ErrDatabaseQuery
// 	}
// 	childPagesDomain := make([]domain.Page, 0)
// 	for _, childPage := range childPages {
// 		childPagesDomain = append(childPagesDomain, *pageutils.TransformPageModelToDomain(childPage, nil))
// 	}

// 	return pageutils.TransformPageModelToDomain(page, childPagesDomain), nil
// }

// func (r *PageRepository) UpdatePageByID(ctx context.Context, pageID string, newPage domain.PageInput) (*domain.Page, *domain.Error) {
// 	var page = model.Page{}

// 	dbErr := r.store.DB().Where("id = ?", pageID).First(&page).Error
// 	if dbErr != nil {
// 		return nil, domain.NewErr("Page not found", domain.BadRequestCode)
// 	}

// 	page.Name = newPage.Name
// 	page.ViewType = newPage.ViewType
// 	page.ParentPagePkid = newPage.ParentPagePkID
// 	page.CoverImage = newPage.CoverImage

// 	dbErr = r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&page).Error

// 	if dbErr != nil {
// 		return nil, domain.ErrDatabaseMutation
// 	}

// 	return pageutils.TransformPageModelToDomain(page, nil), nil
// }

// func (r *PageRepository) ArchivedPageByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error) {
// 	var page = model.Page{}
// 	err := r.store.DB().Where("id = ?", pageID).First(&page).Error
// 	if err != nil {
// 		return nil, domain.NewErr("Page not found", domain.BadRequestCode)
// 	}
// 	now := time.Now()
// 	page.ArchivedAt = &now

// 	err = r.store.DB().Clauses(clause.Returning{}).Select("*").Save(&page).Error
// 	if err != nil {
// 		return nil, domain.ErrDatabaseMutation
// 	}
// 	return pageutils.TransformPageModelToDomain(page, nil), nil
// }

// func (r *PageRepository) GetPagesByNodeID(ctx context.Context, nodeIDs []string) ([]domain.Page, *domain.Error) {
// 	var pages []model.Page
// 	err := r.store.DB().Where("node_id IN ?", nodeIDs).Find(&pages).Error
// 	if err != nil {
// 		return nil, domain.ErrDatabaseQuery
// 	}

// 	var domainPages []domain.Page
// 	for _, page := range pages {
// 		domainPages = append(domainPages, *pageutils.TransformPageModelToDomain(page, nil))
// 	}

// 	return domainPages, nil
// }

// func (r *PageRepository) BulkCreatePages(ctx context.Context, newPagesInput []domain.PageInput) ([]domain.Page, *domain.Error) {
// 	var pages []model.Page
// 	for _, page := range newPagesInput {
// 		nodeID := &page.NodeID
// 		if nodeID == nil {
// 			uuid := uuid.NewString()
// 			nodeID = &uuid
// 		}
// 		pages = append(pages, model.Page{
// 			Name:           page.Name,
// 			SpacePkid:      page.SpacePkID,
// 			ViewType:       page.ViewType,
// 			ParentPagePkid: page.ParentPagePkID,
// 			NodeID:         nodeID,
// 		})
// 	}
// 	err := r.store.DB().Create(&pages).Error
// 	if err != nil {
// 		return nil, domain.ErrDatabaseMutation
// 	}

// 	var domainPages []domain.Page
// 	for _, page := range pages {
// 		domainPages = append(domainPages, *pageutils.TransformPageModelToDomain(page, nil))
// 	}

// 	return domainPages, nil
// }

// func (r *PageRepository) BulkArchivePages(ctx context.Context, pagePkIDs []int64) *domain.Error {
// 	var pages []model.Page
// 	for _, pagePkID := range pagePkIDs {
// 		pages = append(pages, model.Page{
// 			Pkid: pagePkID,
// 		})
// 	}
// 	err := r.store.DB().Model(&pages).Update("archived_at", time.Now()).Error
// 	if err != nil {
// 		return domain.ErrDatabaseMutation
// 	}
// 	return nil
// }
