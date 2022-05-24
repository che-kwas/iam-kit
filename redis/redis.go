// Package redis defines the global redis instance.
package redis // import "github.com/che-kwas/iam-kit/redis"

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	ConfKey = "redis"

	DefaultAddr     = "127.0.0.1:6379"
	DefaultDatabase = 0
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
	log.Printf("NewRedisIns, opts: %+v", opts)

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opts.Addrs,
		Password: opts.Password,
		DB:       opts.Database,
	})

	return rdb, nil
}

func getRedisOpts() (*RedisOptions, error) {
	opts := &RedisOptions{
		Addrs:    []string{DefaultAddr},
		Database: DefaultDatabase,
	}

	if err := viper.UnmarshalKey(ConfKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
