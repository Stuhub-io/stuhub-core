package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	return db, err
}

func Must(dsn string) *gorm.DB {
	db, err := open(dsn)

	if err != nil {
		panic(err)
	}

	return db
}
