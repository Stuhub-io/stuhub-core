package main

import (
	"fmt"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/internal/cache"
	"github.com/Stuhub-io/internal/cache/redis"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/Stuhub-io/logger"
)

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	fmt.Print(cfg)

	logger := logger.NewLogrusLogger()

	postgresDB := postgres.Must(cfg.DBDsn, cfg.Debug, logger)

	redisCache := redis.Must(cfg.RedisUrl)
	cacheStore := cache.NewCacheStore(redisCache)

	dbStore := store.NewDBStore(postgresDB, cacheStore)

	result := dbStore.DB().Exec("ALTER TABLE \"assets\" ALTER COLUMN \"extension\" TYPE VARCHAR(100)")
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Print("Successfully dropped constraint", result)
}
