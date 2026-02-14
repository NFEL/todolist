package cache

import (
	"context"
	"fmt"
	"graph-interview/internal/cfg"
	redis_pkg "graph-interview/pkg/redis"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(cacheCfg cfg.CacheConfig) (*Cache, error) {
	cl, err := redis_pkg.NewClient(
		cacheCfg.Host, cacheCfg.Port, cacheCfg.Password, cacheCfg.DB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating redis client: %w", err)
	}
	return &Cache{Client: cl}, err
}

func (c *Cache) Store(ctx context.Context, key string, value any, duration time.Duration) error {
	cmd := c.Client.Set(ctx, key, value, duration)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	cmd := c.Client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	cmd := c.Client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	cmd := c.Client.Exists(ctx, key)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}
