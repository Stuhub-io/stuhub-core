package postgres

import (
	"context"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/activityutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
	"github.com/Stuhub-io/utils/userutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivityV2Repository struct {
	DB *DB
}

func NewActivityV2Repository(DB *DB) *ActivityV2Repository {
	return &ActivityV2Repository{
		DB: DB,
	}
}

type ActivityResult struct {
	model.Activity
	User *model.User `gorm:"foreignKey:user_pkid"`
}

func (r *ActivityV2Repository) Create(ctx context.Context, input domain.ActivityV2Input) (*domain.ActivityV2, *domain.Error) {
	tx, doneFn := r.DB.NewTransaction()
	activity := &model.Activity{
		UserPkid:   input.UserPkID,
		ActionCode: input.ActionCode.String(),
		Snapshot:   input.Snapshot,
	}

	if err := r.DB.DB().Create(activity).Error; err != nil {
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
	activities := []ActivityResult{}
	if err := buildActivityV2Query(r.DB.DB(), q).Preload("User").Find(&activities).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	domainActivities := make([]domain.ActivityV2, len(activities))

	// FIXME: Query related page Pkids
	for i, activity := range activities {
		domainActivities[i] = *activityutils.TransformActivityV2ModelToDomain(activityutils.ActivityV2ModelToDomainParams{
			ActivityModel:    &activity.Activity,
			RelatedPagePkIDs: nil,
			User:             userutils.TransformUserModelToDomain(activity.User),
		})
	}

	return domainActivities, nil
}

func (r *ActivityV2Repository) One(ctx context.Context, q domain.ActivityV2ListQuery) (*domain.ActivityV2, *domain.Error) {
	activity := &model.Activity{}
	if err := buildActivityV2Query(r.DB.DB(), q).First(activity).Error; err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	// FIXME: Query related page Pkids
	return activityutils.TransformActivityV2ModelToDomain(activityutils.ActivityV2ModelToDomainParams{
		ActivityModel:    activity,
		RelatedPagePkIDs: nil,
	}), nil
}

// NOTE: Cannot remove related related Activity
func (r *ActivityV2Repository) Update(ctx context.Context, activityPkID int64, input domain.ActivityV2Input) (*domain.ActivityV2, *domain.Error) {
	activity := &model.Activity{
		Pkid:       activityPkID,
		UserPkid:   input.UserPkID,
		ActionCode: input.ActionCode.String(),
		Snapshot:   input.Snapshot,
	}

	if err := r.DB.DB().Updates(activity).Error; err != nil {
		return nil, domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}

	// Add related page activity
	RelatedPageActivityList := make([]model.RelatePageActivity, 0, len((input.RelatedPagePkIDs)))
	for _, pagePkID := range input.RelatedPagePkIDs {
		RelatedPageActivityList = append(RelatedPageActivityList, model.RelatePageActivity{
			PagePkid:     pagePkID,
			ActivityPkid: activity.Pkid,
		})
	}

	if err := r.DB.DB().Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&RelatedPageActivityList).Error; err != nil {
		return nil, domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}

	return activityutils.TransformActivityV2ModelToDomain(activityutils.ActivityV2ModelToDomainParams{
		ActivityModel:    activity,
		RelatedPagePkIDs: input.RelatedPagePkIDs,
	}), nil
}

func buildActivityV2Query(tx *gorm.DB, q domain.ActivityV2ListQuery) *gorm.DB {
	query := tx

	ActionCodeStrs := sliceutils.Map(q.ActionCodes, func(code domain.ActionCode) string {
		return code.String()
	})

	if ActionCodeStrs != nil {
		if len(ActionCodeStrs) > 1 {
			query = query.Where("action_code IN ?", ActionCodeStrs)
		} else {
			query = query.Where("action_code = ?", ActionCodeStrs[0])
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

	if q.ForUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	} else {
		query = query.Distinct("activity.pkid", "activity.*")
	}

	return query.Order("created_at DESC")
}
