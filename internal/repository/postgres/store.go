package postgres

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"gorm.io/gorm"
)

// store is implimentation of repository
type store struct {
	database     *gorm.DB
	cacheStore   ports.CacheStore
	shutdownFunc func() error
}

func (s *store) Cache() ports.CacheStore {
	return s.cacheStore
}

// Shutdown close database connection
func (s *store) Shutdown() *domain.Error {
	if s.shutdownFunc != nil {
		if err := s.shutdownFunc(); err != nil {
			return domain.NewErr(err.Error(), domain.InternalServerErrCode)
		}
	}
	return nil
}

// NewTx for database connection
func (s *store) NewTransaction() (repo *ports.Repository, fn ports.IFinallyFunc) {
	newDB := s.database.Begin()

	fn = FinalFunc{db: newDB}
	repo = NewRepo(newDB, s.cacheStore, s.shutdownFunc)

	return repo, fn
}

type FinalFunc struct {
	db *gorm.DB
}

func (fn FinalFunc) Commit() *domain.Error {
	if err := fn.db.Commit().Error; err != nil {
		return domain.NewErr(err.Error(), domain.InternalServerErrCode)
	}
	return nil
}

func (fn FinalFunc) Rollback(err *domain.Error) *domain.Error {
	if rErr := fn.db.Rollback().Error; rErr != nil {
		return domain.NewErr(rErr.Error(), domain.InternalServerErrCode)
	}
	return err
}

// NewStore postgres init by gorm
func NewStore(db *gorm.DB, cache ports.CacheStore, shutdownFunc func() error) ports.DBStore {
	return &store{database: db, cacheStore: cache, shutdownFunc: shutdownFunc}
}
