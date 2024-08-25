package postgres

import (
	"github.com/Stuhub-io/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func open(dsn string, isDebug bool) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logger.L.Fatalf(err, "failed to open database connection")

		return nil, err
	}

	logger.L.Info("database connected")

	if isDebug {
		db.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	return db, nil
}

func Must(dsn string, isDebug bool) *gorm.DB {
	db, err := open(dsn, isDebug)
	if err != nil {
		panic(err)
	}

	return db
}
