// Package redis is the redis instance builder.
package redis // import "github.com/che-kwas/iam-kit/redis"

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKey = "redis"

	defaultAddr     = "127.0.0.1:6379"
	defaultDatabase = 0
)

// RedisOptions defines options for building a redis client.
type RedisOptions struct {
	// Only one addr indicates that it will be a single node client,
	// otherwise it will be a cluster client.
	Addrs    []string
	Password string
	Database int
}

// NewRedisIns creates a redis client.
func NewRedisIns() (redis.UniversalClient, error) {
	opts, err := getRedisOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new redis instance with options: %+v", opts)

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opts.Addrs,
		Password: opts.Password,
		DB:       opts.Database,
	})

	return rdb, nil
}

func getRedisOpts() (*RedisOptions, error) {
	opts := &RedisOptions{
		Addrs:    []string{defaultAddr},
		Database: defaultDatabase,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
