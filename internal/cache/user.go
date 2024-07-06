package cache

import (
	"encoding/json"
	"time"

	"github.com/Stuhub-io/core/domain"
)

func (u *CacheStore) SetUser(user *domain.User, duration time.Duration) error {
	err := u.cache.Set(domain.UserKey(user.PkID), user, duration)

	return err
}

func (u *CacheStore) GetUser(userPkID int64) *domain.User {
	var user domain.User

	data, err := u.cache.Get(domain.UserKey(userPkID))
	if err != nil {
		return nil
	}

	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil
	}

	return &user
}
