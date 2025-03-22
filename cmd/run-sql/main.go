package main

import (
	"fmt"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func open(dsn string, isDebug bool, logger logger.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logger.Fatalf(err, "failed to open database connection")

		return nil, err
	}

	logger.Info("database connected")
	return db, nil
}

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	fmt.Print(cfg)

	logger := logger.NewLogrusLogger()
	db, err := open(cfg.DBDsn, true, logger)
	if err != nil {
		logger.Fatalf(err, "failed to open database connection")
	}

	db.Exec("ALTER TABLE page_access_logs ALTER COLUMN user_pkid DROP NOT NULL;")
	fmt.Print("done")
}
