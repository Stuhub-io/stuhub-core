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
	userPkID int64,
) ([]domain.PageAccessLog, *domain.Error) {
	var result []pageutils.PageAccessLogsResult

	err := r.store.DB().Raw(`
		WITH RECURSIVE parent_pages AS (
			SELECT 
				p.id AS page_id,
				p.pkid AS page_pkid,
				p.name AS page_name,
				p.parent_page_pkid
			FROM pages p

			UNION

			SELECT 
				pp.id AS page_id,
				pp.pkid AS page_pkid,
				pp.name AS page_name,
				pp.parent_page_pkid
			FROM pages pp
			INNER JOIN parent_pages p ON pp.pkid = p.parent_page_pkid
		)
		SELECT 
			pl.pkid, 
			p.pkid AS page_pkid,
			p.id AS page_id,
			p.name AS page_name,
			pl.action, 
			CASE 
				WHEN d.page_pkid IS NOT NULL THEN 'document'
				WHEN a.page_pkid IS NOT NULL THEN 'asset'
				ELSE 'none'
			END AS view_type, 
			u.first_name, 
			u.last_name, 
			u.email, 
    		u.avatar, 
			pl.last_accessed,
			ARRAY(
				SELECT json_build_object(
					'id', parent_pages.page_id,
					'pkid', parent_pages.page_pkid,
					'name', parent_pages.page_name
				)
				FROM parent_pages
				WHERE parent_pages.page_pkid IN (
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
		WHERE pl.user_pkid = ?`, userPkID).Scan(&result).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	var accessLogs []domain.PageAccessLog
	for _, row := range result {
		accessLogs = append(accessLogs, pageutils.TransformPageAccessLogsResultToDomain(row))
	}

	return accessLogs, nil
}
