package postgres

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func open(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

func Must(dsn string) *bun.DB {
	db := open(dsn)

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("could not connect to DB")
	}

	fmt.Println("Connected to DB successfully!")

	return db
}
