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

func NewCache(cfg cfg.CacheConfig) (*Cache, error) {
	cl, err := redis_pkg.NewClient(
		cfg.Host, cfg.Port, cfg.Password, cfg.DB,
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

func (c *Cache) Get(ctx context.Context, key string, value any, duration time.Duration) error {
	cmd := c.Client.Set(ctx, key, value, duration)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}
