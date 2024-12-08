package postgres

import (
	"context"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/pageutils"
	"gorm.io/gorm/clause"
)

func (r *PageRepository) CreatePublicToken(ctx context.Context, pagePkID int64) (*domain.PagePublicToken, *domain.Error) {
	var page model.Page
	if dbErr := r.store.DB().Where("pkid = ?", pagePkID).First(&page).Error; dbErr != nil {
		return nil, domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	newToken := model.PublicToken{
		PagePkid: page.Pkid,
	}
	err := r.store.DB().Create(&newToken).Error
	if err != nil {
		return nil, domain.NewErr(err.Error(), domain.BadRequestCode)
	}

	return pageutils.TransformPagePublicTokenModelToDomain(newToken), nil
}

func (r *PageRepository) ArchiveAllPublicToken(ctx context.Context, pagePkID int64) *domain.Error {

	now := time.Now()

	if dbErr := r.store.DB().Clauses(clause.Returning{}).
		Model(&model.PublicToken{}).
		Where("page_pkid = ?", pagePkID).
		Select("ArchivedAt").
		Updates(model.PublicToken{
			ArchivedAt: &now,
		}).Error; dbErr != nil {
		return domain.NewErr(dbErr.Error(), domain.BadRequestCode)
	}

	return nil
}

func (p *PageRepository) GetPublicTokenByID(ctx context.Context, publicTokenID string) (*domain.PagePublicToken, *domain.Error) {
	var token model.PublicToken
	if dbErr := p.store.DB().Where("id = ?", publicTokenID).First(&token).Error; dbErr != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return pageutils.TransformPagePublicTokenModelToDomain(token), nil
}
