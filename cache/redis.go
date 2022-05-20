// Package cache is a cache instance builder.
package cache // import "github.com/che-kwas/iam-kit/cache"

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const ConfigCacheKey = "redis"

// CacheBuilder defines options for building a cache instance.
type CacheBuilder struct {
	Addrs    []string
	Password string
	Database int

	err error
}

// NewCacheBuilder is used to build a cache instance.
func NewCacheBuilder() *CacheBuilder {
	b := &CacheBuilder{}
	b.err = viper.UnmarshalKey(ConfigCacheKey, b)

	return b
}

// Build builds a redis cache instance.
func (b *CacheBuilder) Build() (*redis.UniversalClient, error) {
	if b.err != nil {
		return nil, b.err
	}

	// rdb is *redis.ClusterClient.
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    b.Addrs,
		Password: b.Password,
		DB:       b.Database,
	})

	return &rdb, nil
}
