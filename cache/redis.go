// Package cache is a cache instance builder.
package cache // import "github.com/che-kwas/iam-kit/cache"

import (
	"github.com/che-kwas/iam-kit/config"
	"github.com/go-redis/redis/v8"
)

// NewRedisIns builds a redis cache instance.
func NewRedisIns(opts *config.RedisOptions) (*redis.UniversalClient, error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opts.Addrs,
		Password: opts.Password,
		DB:       opts.Database,
	})

	return &rdb, nil
}
