package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reading-microservices/user-service/models"
	"time"
)

type UserCache interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	SetUser(ctx context.Context, user *models.User, expiration time.Duration) error
	InvalidateUser(ctx context.Context, userID string) error
	GetSession(ctx context.Context, token string) (*models.UserSession, error)
	SetSession(ctx context.Context, session *models.UserSession, expiration time.Duration) error
	IncrementLoginAttempts(ctx context.Context, username string) (int64, error)
	GetLoginAttempts(ctx context.Context, username string) (int64, error)
	ResetLoginAttempts(ctx context.Context, username string) error
}

type redisUserCache struct {
	client *redis.Client
	prefix string
}

func NewRedisUserCache(client *redis.Client) UserCache {
	return &redisUserCache{
		client: client,
		prefix: "user_service:",
	}
}

func (c *redisUserCache) GetUser(ctx context.Context, userID string) (*models.User, error) {
	key := fmt.Sprintf("%suser:%s", c.prefix, userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *redisUserCache) SetUser(ctx context.Context, user *models.User, expiration time.Duration) error {
	key := fmt.Sprintf("%suser:%s", c.prefix, user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *redisUserCache) InvalidateUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("%suser:%s", c.prefix, userID)
	return c.client.Del(ctx, key).Err()
}

func (c *redisUserCache) GetSession(ctx context.Context, token string) (*models.UserSession, error) {
	key := fmt.Sprintf("%ssession:%s", c.prefix, token)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var session models.UserSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (c *redisUserCache) SetSession(ctx context.Context, session *models.UserSession, expiration time.Duration) error {
	key := fmt.Sprintf("%ssession:%s", c.prefix, session.SessionToken)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *redisUserCache) IncrementLoginAttempts(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%slogin_attempts:%s", c.prefix, username)
	return c.client.Incr(ctx, key).Result()
}

func (c *redisUserCache) GetLoginAttempts(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("%slogin_attempts:%s", c.prefix, username)
	return c.client.Get(ctx, key).Int64()
}

func (c *redisUserCache) ResetLoginAttempts(ctx context.Context, username string) error {
	key := fmt.Sprintf("%slogin_attempts:%s", c.prefix, username)
	return c.client.Del(ctx, key).Err()
}
