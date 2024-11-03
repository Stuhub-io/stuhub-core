package main

import (
	"fmt"

	"github.com/Stuhub-io/config"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/google/uuid"
)

func main() {

	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())
	postgresDB := postgres.Must(cfg.DBDsn, cfg.Debug)
	dbStore := store.NewDBStore(postgresDB, nil)

	fmt.Println("Successfully connected to the database")

	var missingNodeIDPages []model.Page
	err := dbStore.DB().Where("node_id IS NULL").Find(&missingNodeIDPages).Error
	if err != nil {
		panic(err)
	}

	for i := range missingNodeIDPages {
		newUUID := uuid.NewString()
		missingNodeIDPages[i].NodeID = &newUUID
	}
	rerr := dbStore.DB().Save(&missingNodeIDPages).Error
	if rerr != nil {
		panic(rerr)
	}
	fmt.Println("Successfully generated node IDs for missing pages")
}
