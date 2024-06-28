package store

import (
	"errors"

	"gorm.io/gorm"
)

type TxEndFunc func(error) error
type DBStore interface {
	DB() *gorm.DB
	NewTransaction() (DBStore, TxEndFunc)
	SetNewDB(*gorm.DB)
}

func NewDbStore(db *gorm.DB) DBStore {
	return &dbStore{
		Database: db,
	}
}

type dbStore struct {
	Database *gorm.DB
}

func (d *dbStore) DB() *gorm.DB {
	return d.Database
}

func (d *dbStore) NewTransaction() (DBStore, TxEndFunc) {
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

	return &dbStore{Database: newDB}, finallyFn
}

func (d *dbStore) SetNewDB(db *gorm.DB) {
	d.Database = db
}
