package main

import (
	"fmt"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/internal/cache"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/Stuhub-io/logger"
)

type TempCache struct{}

func (TempCache) Set(key string, value any, duration time.Duration) error { return nil }
func (TempCache) Get(key string) (string, error)                          { return "", nil }
func (TempCache) Delete(key string) error                                 { return nil }

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	fmt.Print(cfg)

	logger := logger.NewLogrusLogger()

	postgresDB := postgres.Must(cfg.DBDsn, cfg.Debug, logger)

	// redisCache := redis.Must(cfg.RedisUrl)
	tempCache := TempCache{}
	cacheStore := cache.NewCacheStore(tempCache)

	dbStore := store.NewDBStore(postgresDB, cacheStore)

	// modify the sql
	result := dbStore.DB().Exec(`
			ALTER TABLE pages
			DROP CONSTRAINT IF EXISTS fk_page_author;

			-- Add the foreign key constraint with ON DELETE SET NULL
			ALTER TABLE pages
			ADD CONSTRAINT fk_page_author
			FOREIGN KEY (author_pkid) REFERENCES users(pkid)
			ON DELETE SET NULL;
	`)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Print("Successfully dropped constraint", result)
}
