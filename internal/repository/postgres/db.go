package postgres

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func open(dsn string, isDebug bool, logger logger.Logger) (*gorm.DB, error) {
	conn := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(conn, &gorm.Config{})

	if err != nil {
		logger.Fatalf(err, "failed to open database connection")

		return nil, err
	}

	logger.Info("database connected")

	if isDebug {
		db.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	return db, nil
}

func Must(dsn string, isDebug bool, l logger.Logger) *gorm.DB {
	db, err := open(dsn, isDebug, l)
	if err != nil {
		panic(err)
	}

	return db
}

type DB struct {
	Database *gorm.DB // Primary DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		Database: db,
	}
}

func (d *DB) DB() *gorm.DB {
	return d.Database
}

type TxEndFunc func(error) *domain.Error

func (d *DB) NewTransaction() (*DB, TxEndFunc) {
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

	return &DB{Database: newDB}, finallyFn
}

func (d *DB) SetNewDB(db *gorm.DB) {
	d.Database = db
}
