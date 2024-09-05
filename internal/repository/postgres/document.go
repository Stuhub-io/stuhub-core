package postgres

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/docutils"
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

func (r *DocRepository) CreateDocument(ctx context.Context, pagePkID int64, JsonContent string) (*domain.Document, *domain.Error) {
	newDoc := model.Document{
		PagePkid:    pagePkID,
		Content:     "",
		JSONContent: &JsonContent,
	}
	err := r.store.DB().Create(&newDoc).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	return docutils.TransformDocModalToDomain(newDoc), nil
}

func (r *DocRepository) GetDocumentByPagePkID(ctx context.Context, pagePkID int64) (*domain.Document, *domain.Error) {
	var doc model.Document
	err := r.store.DB().Where("page_pkid = ?", pagePkID).Order("created_at desc").First(&doc).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	return docutils.TransformDocModalToDomain(doc), nil
}

func (r *DocRepository) GetDocumentByPkID(ctx context.Context, pkID int64) (*domain.Document, *domain.Error) {
	var doc model.Document
	err := r.store.DB().Where("pkid = ?", pkID).Order("created_at desc").First(&doc).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	return docutils.TransformDocModalToDomain(doc), nil
}

func (r *DocRepository) UpdateDocument(ctx context.Context, pkID int64, jsonContent string) (*domain.Document, *domain.Error) {
	var doc model.Document
	err := r.store.DB().Where("pkid = ?", pkID).Order("created_at desc").First(&doc).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	r.store.DB().Model(&doc).Updates(map[string]interface{}{"json_content": jsonContent})

	return docutils.TransformDocModalToDomain(doc), nil
}
