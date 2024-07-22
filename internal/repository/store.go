package store

import (
	"errors"

	"github.com/Stuhub-io/core/ports"
	"gorm.io/gorm"
)

type TxEndFunc func(error) error

type DBStore struct {
	Database   *gorm.DB
	CacheStore ports.CacheStore
}

func NewDBStore(db *gorm.DB, cacheStore ports.CacheStore) *DBStore {
	return &DBStore{
		Database:   db,
		CacheStore: cacheStore,
	}
}

func (d *DBStore) DB() *gorm.DB {
	return d.Database
}

func (d *DBStore) Cache() ports.CacheStore {
	return d.CacheStore
}

func (d *DBStore) NewTransaction() (*DBStore, TxEndFunc) {
	newDB := d.Database.Begin()

	finallyFn := func(err error) error {
		if err != nil {
			nErr := newDB.Rollback().Error
			if nErr != nil {
				return errors.New(nErr.Error())
			}

			return err
		}

		cErr := newDB.Commit().Error
		if cErr != nil {
			return errors.New(cErr.Error())
		}

		return nil
	}

	return &DBStore{Database: newDB}, finallyFn
}

func (d *DBStore) SetNewDB(db *gorm.DB) {
	d.Database = db
}
