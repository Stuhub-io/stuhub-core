package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/internal/cache"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
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

	// Migrate All Paths From IDs Into PkIDs
	var allPages []model.Page
	dbStore.DB().Find(&allPages)

	pageIDMap := make(map[string]model.Page)
	for _, page := range allPages {
		pageIDMap[page.ID] = page
	}

	PathToPkIDPath := func(path string) string {
		if len(path) == 0 {
			return ""
		}
		ids := strings.Split(path, "/")
		pkIDs := make([]string, 0, len(ids))
		for _, id := range ids {
			pkIdStr := strconv.FormatInt(pageIDMap[id].Pkid, 10)
			pkIDs = append(pkIDs, pkIdStr)
		}
		newPath := strings.Join(pkIDs, "/")
		return newPath
	}

	updatePages := make([]model.Page, 0, len(allPages))

	for _, page := range allPages {
		page.Path = PathToPkIDPath(page.Path)
		updatePages = append(updatePages, page)
		fmt.Println(page.Path)
	}

	dbStore.DB().Save(&updatePages)

	fmt.Print("done")
}
