package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	commonutils "github.com/Stuhub-io/utils"
	"gorm.io/gorm/clause"
)

func (r *PageRepository) StarPage(ctx context.Context, input domain.StarPageInput) (*domain.PageStar, *domain.Error) {
	starPage := model.PageStar{
		UserPkid: input.ActorUserPkID,
		PagePkid: input.PagePkID,
		Order:    commonutils.CurTimestampAsFloat64(), // Order is used to sort the starred pages
	}
	if err := r.store.DB().Clauses(clause.Returning{}).Create(&starPage).Error; err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return nil, nil
}

func (r *PageRepository) UnstarPage(ctx context.Context, input domain.StarPageInput) *domain.Error {
	q := r.store.DB().Where("user_pkid = ? AND page_pkid = ?", input.ActorUserPkID, input.PagePkID)
	if err := q.Delete(&model.PageStar{}).Error; err != nil {
		return domain.ErrDatabaseMutation
	}
	return nil
}
