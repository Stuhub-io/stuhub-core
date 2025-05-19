package ports

import (
	"github.com/Stuhub-io/core/domain"
	"gorm.io/gorm"
)

type TxEndFunc func(error) *domain.Error

type DBStore interface {
	DB() *gorm.DB
	Cache() CacheStore
	NewTransaction() (DBStore, TxEndFunc)
	SetNewDB(db *gorm.DB)
}
