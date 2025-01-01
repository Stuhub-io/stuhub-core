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
		ParentName   *string `gorm:"column:name"`
		LastAccessed time.Time
	}

	err := r.store.DB().Raw(`
		SELECT 
			pl.pkid, 
			pl.action, 
			p.name,
   			CASE 
				WHEN d.page_pkid IS NOT NULL THEN 'document'
				WHEN a.page_pkid IS NOT NULL THEN 'asset'
				ELSE 'none'
			END AS view_type, 
			u.email, 
			u.avatar, 
			pp.name as folder_name,
			pl.last_accessed
		FROM page_access_logs pl 
		LEFT JOIN pages p ON p.pkid = pl.page_pkid 
		LEFT JOIN pages pp ON p.parent_page_pkid = pp.pkid
		LEFT JOIN users u ON u.pkid = p.author_pkid
		LEFT JOIN documents d ON d.page_pkid = p.pkid
		LEFT JOIN assets a ON a.page_pkid = p.pkid
		WHERE pl.user_pkid= ?`, userPkID).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	fmt.Print("hello ", result)
	return result, nil
}
