// Package redis defines the global redis instance.
package redis // import "github.com/che-kwas/iam-kit/redis"

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	ConfKey = "redis"

	DefaultAddr     = "127.0.0.1:6379"
	DefaultDatabase = 0
)

var (
	rdb  redis.UniversalClient
	once sync.Once
)

// RedisOptions defines options for building a redis client.
type RedisOptions struct {
	// Only one addr indicates that it will be a single node client,
	// otherwise it will be a cluster client.
	Addrs    []string
	Password string
	Database int
}

// GetRedisIns returns a redis client.
func GetRedisIns() (redis.UniversalClient, error) {
	if rdb != nil {
		return rdb, nil
	}

	var err error
	once.Do(func() { rdb, err = newRedisIns() })

	return rdb, err

}

func newRedisIns() (redis.UniversalClient, error) {
	opts, err := getRedisOpts()
	if err != nil {
		return nil, err
	}

	rdb = redis.NewUniversalClient(&redis.UniversalOptions{
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
