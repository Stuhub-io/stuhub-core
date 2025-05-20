package postgres

import (
	"github.com/Stuhub-io/core/ports"
	"gorm.io/gorm"
)

func NewRepo(DB *gorm.DB, Cache ports.CacheStore, ShutdownFunc func() error) *ports.Repository {
	db := NewDB(DB)
	usrRepo := NewUserRepository(db)
	return &ports.Repository{
		Store:              NewStore(db.DB(), Cache, ShutdownFunc),
		User:               usrRepo,
		Organization:       NewOrganizationRepository(db, usrRepo),
		Page:               NewPageRepository(db),
		PageAccessLog:      NewPageAccessLogRepository(db),
		OrganizationInvite: NewOrganizationInvitesRepository(db),
		ActivityV2:         NewActivityV2Repository(db),
	}
}
