package cache

import "github.com/Stuhub-io/core/ports"

type CacheStore struct {
	cache ports.Cache
}

func NewCacheStore(cache ports.Cache) *CacheStore {
	return &CacheStore{
		cache: cache,
	}
}
