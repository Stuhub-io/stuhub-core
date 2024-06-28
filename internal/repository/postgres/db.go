package postgres

import (
	"database/sql"

	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func open(dsn string, isDebug bool) (store.DBStore, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.L.Fatalf(err, "failed to open database connection")
		return nil, err
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: conn}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
		})

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
