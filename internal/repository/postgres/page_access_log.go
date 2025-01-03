package postgres

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/utils/pageutils"
)

type PageAccessLogRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewPageAccessLogRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewPageAccessLogRepository(params NewPageAccessLogRepositoryParams) *PageAccessLogRepository {
	return &PageAccessLogRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *PageAccessLogRepository) GetByUserPKID(
	ctx context.Context,
	query domain.OffsetBasedPagination,
	userPkID int64,
) ([]domain.PageAccessLog, *domain.Error) {
	var result []pageutils.PageAccessLogsResult

	err := r.store.DB().Raw(`
		SELECT 
			pl.pkid, 
			p.pkid AS page_pkid,
			p.id AS page_id,
			p.name AS page_name,
			p.general_role AS page_general_role,
			pl.action, 
			CASE 
				WHEN d.page_pkid IS NOT NULL THEN 'document'
				WHEN a.page_pkid IS NOT NULL THEN 'asset'
				ELSE 'none'
			END AS view_type, 
			u.pkid as author_pkid, 
			u.first_name as author_first_name, 
			u.last_name as author_last_name, 
			u.email as author_email, 
    		u.avatar as author_avatar, 
			pl.last_accessed,
			ARRAY(
				SELECT json_build_object(
					'id', pages.id,
					'pkid', pages.pkid,
					'name', pages.name,
					'author_pkid', pages.author_pkid,
					'general_role', pages.general_role
				)
				FROM pages
				WHERE pages.pkid IN (
					WITH RECURSIVE page_hierarchy AS (
						SELECT pkid, parent_page_pkid
						FROM pages
						WHERE pkid = pl.page_pkid
						
						UNION
						
						SELECT p2.pkid, p2.parent_page_pkid
						FROM pages p2
						INNER JOIN page_hierarchy ph ON p2.pkid = ph.parent_page_pkid
					)
					SELECT pkid FROM page_hierarchy WHERE pkid != p.pkid
				)
			) AS parent_pages
		FROM page_access_logs pl 
		LEFT JOIN pages p ON p.pkid = pl.page_pkid 
		LEFT JOIN users u ON u.pkid = p.author_pkid
		LEFT JOIN documents d ON d.page_pkid = p.pkid
		LEFT JOIN assets a ON a.page_pkid = p.pkid
		WHERE pl.user_pkid = ? LIMIT ? OFFSET ?`, userPkID, query.Limit, query.Offset).Scan(&result).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var accessLogs []domain.PageAccessLog
	for _, row := range result {
		accessLogs = append(accessLogs, pageutils.TransformPageAccessLogsResultToDomain(row))
	}

	return accessLogs, nil
}