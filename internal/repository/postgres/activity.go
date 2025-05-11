package postgres

import (
	"context"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/activityutils"
	"gorm.io/gorm"
)

type ActivityV2Repository struct {
	cfg   config.Config
	store *store.DBStore
}

type ActivityV2RepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewActivityV2Repository(params ActivityV2RepositoryParams) *ActivityV2Repository {
	return &ActivityV2Repository{
		cfg:   params.Cfg,
		store: params.Store,
	}
}

func (r *ActivityV2Repository) Create(ctx context.Context, input domain.ActivityV2Input) (*domain.ActivityV2, *domain.Error) {
	tx, doneFn := r.store.NewTransaction()
	activity := &model.Activity{
		UserPkid:   input.UserPkID,
		ActionCode: input.ActionCode.String(),
		Snapshot:   input.Snapshot,
	}

	if err := tx.DB().Create(activity).Error; err != nil {
		return nil, doneFn(err)
	}

	RelatedPageActivityList := make([]model.RelatePageActivity, 0, len((input.RelatedPagePkIDs)))
	for _, pagePkID := range input.RelatedPagePkIDs {
		RelatedPageActivityList = append(RelatedPageActivityList, model.RelatePageActivity{
			PagePkid:     pagePkID,
			ActivityPkid: activity.Pkid,
		})
	}

	if err := tx.DB().Create(&RelatedPageActivityList).Error; err != nil {
		return nil, doneFn(err)
	}

	return activityutils.TransformActivityV2ModelToDomain(activityutils.ActivityV2ModelToDomainParams{
		ActivityModel:    activity,
		RelatedPagePkIDs: input.RelatedPagePkIDs,
	}), doneFn(nil)
}

func (r *ActivityV2Repository) List(ctx context.Context, q domain.ActivityV2ListQuery) ([]domain.ActivityV2, *domain.Error) {
	activities := []model.Activity{}
	if err := buildActivityV2Query(r.store.DB(), q).Find(&activities).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	domainActivities := make([]domain.ActivityV2, len(activities))
	for i, activity := range activities {
		domainActivities[i] = *activityutils.TransformActivityV2ModelToDomain(activityutils.ActivityV2ModelToDomainParams{
			ActivityModel:    &activity,
			RelatedPagePkIDs: nil,
		})
	}

	return domainActivities, nil
}

func buildActivityV2Query(tx *gorm.DB, q domain.ActivityV2ListQuery) *gorm.DB {
	query := tx.Distinct("activity.pkid", "activity.*")

	if q.ActionCodes != nil {
		if len(q.ActionCodes) > 1 {
			query = query.Where("action_code IN ?", q.ActionCodes)
		} else {
			query = query.Where("action_code = ?", q.ActionCodes[0])
		}
	}

	if q.EndTime != nil {
		query = query.Where("created_at <= ?", q.EndTime.Format(time.RFC3339))
	}

	if q.UserPkIDs != nil {
		if len(q.UserPkIDs) > 1 {
			query = query.Where("user_pkid IN ?", q.UserPkIDs)
		} else {
			query = query.Where("user_pkid = ?", q.UserPkIDs[0])
		}
	}

	if q.RelatedPagePkIDs != nil {
		query = query.Joins("JOIN relate_page_activity ON relate_page_activity.activity_pkid = activity.pkid")
		if len(q.RelatedPagePkIDs) > 1 {
			query = query.Where("relate_page_activity.page_pkid IN ?", q.RelatedPagePkIDs)
		} else {
			query = query.Where("relate_page_activity.page_pkid = ?", q.RelatedPagePkIDs[0])
		}
	}

	if q.Limit != nil {
		query = query.Limit(*q.Limit)
	}

	return query
}
