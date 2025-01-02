package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Stuhub-io/config"
	store "github.com/Stuhub-io/internal/repository"
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

func (r *PageAccessLogRepository) GetByUserPKID(ctx context.Context, userPkID int64) (any, error) {
	var result []struct {
		Pkid         int64
		Action       string
		Name         string
		ViewType     string
		Email        string
		Avatar       string
		LastAccessed time.Time
		Pages        string
	}
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
			pl.action, 
			p.name AS page_name,
			CASE 
				WHEN d.page_pkid IS NOT NULL THEN 'document'
				WHEN a.page_pkid IS NOT NULL THEN 'asset'
				ELSE 'none'
			END AS view_type, 
			u.email, 
			u.avatar, 
			pl.last_accessed,
			ARRAY(
				SELECT json_build_object(
					'id', parent_pages.page_id,
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
					SELECT pkid FROM page_hierarchy
				)
			) AS pages
		FROM page_access_logs pl 
		LEFT JOIN pages p ON p.pkid = pl.page_pkid 
		LEFT JOIN users u ON u.pkid = p.author_pkid
		LEFT JOIN documents d ON d.page_pkid = p.pkid
		LEFT JOIN assets a ON a.page_pkid = p.pkid
		WHERE pl.user_pkid = ?`, userPkID).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	fmt.Print("hello ", result)
	return result, nil
}
