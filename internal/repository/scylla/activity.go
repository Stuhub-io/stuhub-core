package scylla

import (
	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
)

type ActivityRepository struct {
	cfg   config.Config
	store *store.DBStore
}

type ActivityRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewActivityRepository(params ActivityRepositoryParams) *ActivityRepository {
	return &ActivityRepository{
		cfg:   params.Cfg,
		store: params.Store,
	}
}

func (r *ActivityRepository) List() ([]domain.Activity, *domain.Error) {
	r.store.LogDB()
	return nil, nil
}
