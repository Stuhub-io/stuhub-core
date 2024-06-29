package postgres

import (
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func open(dsn string, isDebug bool) (store.DBStore, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		logger.L.Fatalf(err, "failed to open database connection")

		return nil, err
	}

	logger.L.Info("database connected")

	if isDebug {
		db.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	return store.NewDbStore(db), nil
}

func NewStore(dsn string, isDebug bool) (store.DBStore, error) {
	return open(dsn, isDebug)
}
