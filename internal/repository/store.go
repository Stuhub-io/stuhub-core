package store

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/gocql/gocql"
	"gorm.io/gorm"
)

type DBStore struct {
	Database    *gorm.DB         // Primary DB
	LogDatabase *gocql.Session   // Secondary DB
	CacheStore  ports.CacheStore // Cache store
}

func NewDBStore(db *gorm.DB, cacheStore ports.CacheStore, logDb *gocql.Session) *DBStore {
	return &DBStore{
		Database:    db,
		CacheStore:  cacheStore,
		LogDatabase: logDb,
	}
}

func (d *DBStore) DB() *gorm.DB {
	return d.Database
}

func (d *DBStore) Cache() ports.CacheStore {
	return d.CacheStore
}

func (d *DBStore) NewTransaction() (ports.DBStore, ports.TxEndFunc) {
	newDB := d.Database.Begin()

	finallyFn := func(err error) *domain.Error {
		if err != nil {
			nErr := newDB.Rollback().Error
			if nErr != nil {
				return domain.NewErr(err.Error(), domain.InternalServerErrCode)
			}

			return domain.ErrInternalServerError
		}

		cErr := newDB.Commit().Error
		if cErr != nil {
			return domain.NewErr(cErr.Error(), domain.InternalServerErrCode)
		}

		return nil
	}

	return &DBStore{Database: newDB}, finallyFn
}

func (d *DBStore) SetNewDB(db *gorm.DB) {
	d.Database = db
}

func (d *DBStore) LogDB() *gocql.Session {
	return d.LogDatabase
}
