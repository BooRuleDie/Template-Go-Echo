package user

import (
	"context"
	"encoding/json"

	keys "go-echo-template/internal/cache"
	"go-echo-template/internal/storage/user/sqlc"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Get(ctx context.Context, userID int64) (*sqlc.User, error)
	Set(ctx context.Context, user *sqlc.User) error
	Delete(ctx context.Context, userID int64) error
}

type cache struct {
	rc *redis.Client
}

func NewUserCache(rc *redis.Client) UserCache {
	return &cache{rc: rc}
}

func (c *cache) Get(ctx context.Context, userID int64) (*sqlc.User, error) {
	cacheKey := keys.GetUserKey(userID)

	data, err := c.rc.Get(ctx, cacheKey.Name).Result()
	if err != nil {
		if err == redis.Nil {
			// not an error actually, no data found
			return nil, nil
		} else {
			// actual error
			return nil, err
		}
	}

	user := new(sqlc.User)
	if data != "" {
		if err := json.Unmarshal([]byte(data), user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (c *cache) Set(ctx context.Context, user *sqlc.User) error {
	cacheKey := keys.GetUserKey(user.ID)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.rc.SetEx(ctx, cacheKey.Name, jsonData, cacheKey.TTL).Err()
}

func (c *cache) Delete(ctx context.Context, userID int64) error {
	cacheKey := keys.GetUserKey(userID)
	return c.rc.Del(ctx, cacheKey.Name).Err()
}
