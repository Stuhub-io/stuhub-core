package postgres

import (
	"context"
	"fmt"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
)

type PageRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewPageRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewPageRepository(params NewPageRepositoryParams) ports.PageRepository {
	return &PageRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *PageRepository) CreatePage(ctx context.Context, spacePkID int64, name string, viewType domain.PageViewType, ParentPagePkID *int64) (*domain.Page, *domain.Error) {
	fmt.Print("\n\n\n", ParentPagePkID, ParentPagePkID == nil, "\n\n\n")
	// var newPage model.Page
	newPage := model.Page{
		Name:           name,
		SpacePkid:      spacePkID,
		ViewType:       viewType.String(),
		ParentPagePkid: ParentPagePkID,
	}
	err := r.store.DB().Create(&newPage).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	return pageutils.MapPageModelToDomain(newPage), nil
}

func (r *PageRepository) GetPagesBySpacePkID(ctx context.Context, spacePkID int64) ([]domain.Page, *domain.Error) {
	var pages []model.Page
	err := r.store.DB().Where("space_pkid = ?", spacePkID).Find(&pages).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var domainPages []domain.Page
	for _, page := range pages {
		domainPages = append(domainPages, *pageutils.MapPageModelToDomain((page)))
	}

	return domainPages, nil
}