package scylla

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (r *ActivityRepository) List(ctx context.Context, q domain.ActivityListQuery) ([]domain.Activity, *domain.Error) {

	queryStr, args := buildActivityQuery(q)
	iter := r.store.LogDB().Query(queryStr, args...).PageSize(q.Limit).Iter()

	activities := make([]domain.Activity, 0, iter.NumRows())

	curActivity := domain.Activity{}
	var createdTime time.Time

	for iter.Scan(
		&curActivity.ActorPkID,
		&curActivity.PagePkID,
		&curActivity.OrgPkID,
		&curActivity.ActionCode,
		&curActivity.Label,
		&curActivity.MetaData,
		&createdTime,
	) {
		curActivity.CreatedAt = createdTime.Format(time.RFC3339)
		activities = append(activities, curActivity)
	}

	if err := iter.Close(); err != nil {
		return nil, domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}

	return activities, nil
}

func (r *ActivityRepository) Create(ctx context.Context, input domain.ActivityInput) (*domain.Activity, *domain.Error) {
	// Get the current time for created_at
	createdAt := time.Now().UTC()
	if err := r.store.LogDB().Query(
		`INSERT INTO activity (org_pkid, actor_pkid, page_pkid, action_code, label, metadata, created_at) VALUES (?,?, ?, ?, ?, ?, ?)`,
		input.OrgPkID,
		input.ActorPkID,
		input.PagePkID,
		input.ActionCode,
		input.Label,
		input.MetaData,
		createdAt,
	).Exec(); err != nil {
		return nil, domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}

	return &domain.Activity{
		ActorPkID:  input.ActorPkID,
		PagePkID:   input.PagePkID,
		ActionCode: input.ActionCode,
		Label:      input.Label,
		MetaData:   input.MetaData,
		OrgPkID:    input.OrgPkID,
		CreatedAt:  createdAt.Format(time.RFC3339),
	}, nil
}

func buildActivityQuery(query domain.ActivityListQuery) (string, []interface{}) {
	baseQuery := `SELECT actor_pkid, page_pkid, org_pkid, action_code, label, metadata, created_at FROM activity`

	var conditions []string
	var args []interface{}

	// Filter By `ACTION_CODE`
	if len(query.ActionCodes) > 0 {
		placeholders := make([]string, len(query.ActionCodes))
		for i, code := range query.ActionCodes {
			placeholders[i] = fmt.Sprintf("?")
			args = append(args, string(code))
		}
		conditions = append(conditions, fmt.Sprintf("action_code IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add filter for ActorPkIDs if provided
	if len(query.ActorPkIDs) > 0 {
		placeholders := make([]string, len(query.ActorPkIDs))
		for i, id := range query.ActorPkIDs {
			placeholders[i] = fmt.Sprintf("?")
			args = append(args, id)
		}
		conditions = append(conditions, fmt.Sprintf("actor_pkid IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add filter for PagePkIDs if provided
	if len(query.PagePkIDs) > 0 {
		placeholders := make([]string, len(query.PagePkIDs))
		for i, id := range query.PagePkIDs {
			placeholders[i] = fmt.Sprintf("?")
			args = append(args, id)
		}
		conditions = append(conditions, fmt.Sprintf("page_pkid IN (%s)", strings.Join(placeholders, ", ")))
	}

	if query.EndTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < ?"))
		args = append(args, *query.EndTime)
	}

	if query.StartTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at > ?"))
		args = append(args, *query.StartTime)
	}

	// Construct the final query
	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add semicolon at the end of the query
	finalQuery += " ALLOW FILTERING;"

	return finalQuery, args
}
