package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthData struct {
	UserId   uuid.UUID
	Email    string
	FullName string
	Role     string
	IsActive bool
}

type AuthCache struct {
	Redis *redis.Client
	Ctx   context.Context
}

func NewAuthCache(redis *redis.Client) *AuthCache {
	return &AuthCache{
		Redis: redis,
		Ctx:   context.Background(),
	}
}

func (authCache *AuthCache) GetAuthCache(key string) (*AuthData, error) {
	val, err := authCache.Redis.Get(authCache.Ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	authData := &AuthData{}
	if err := json.Unmarshal([]byte(val), authData); err != nil {
		return nil, err
	}
	return authData, nil
}

func (authCache *AuthCache) SetAuthCache(key string, value *AuthData, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = authCache.Redis.Set(authCache.Ctx, key, string(b), ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (authCache *AuthCache) SetRefreshTokenCache(key string, refersh string, ttl time.Duration) error {
	err := authCache.Redis.Set(authCache.Ctx, key, refersh, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (authCache *AuthCache) GetRefreshTokenCache(key string) (string, error) {
	value, err := authCache.Redis.Get(authCache.Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (authCache *AuthCache) DeleteAuthCache(key string) error {
	err := authCache.Redis.Del(authCache.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
