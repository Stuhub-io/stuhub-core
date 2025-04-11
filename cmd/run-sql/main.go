package main

import (
	"fmt"
	"strconv"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/logger"
	"github.com/gocql/gocql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string, isDebug bool, logger logger.Logger) (*gorm.DB, error) {
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

func OpenScylla(hosts []string, port string, keyspace string, isDebug bool, logger logger.Logger) *gocql.Session {
	nPort, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatalf(err, "failed to convert port to int: %v")
		panic(err)
	}
	cluster := gocql.NewCluster(hosts...)
	cluster.Port = nPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		logger.Fatalf(err, "failed to create Cassandra session: %v")
		panic(err)
	}

	if isDebug {
		session.SetConsistency(gocql.Quorum)
	}

	logger.Info("SCylla connected")
	return session
}

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	logger := logger.NewLogrusLogger()

	db, err := Open(cfg.DBDsn, true, logger)
	if err != nil {
		logger.Fatalf(err, "failed to open database connection")
	}

	db.Exec("")
	fmt.Print("done")
}
