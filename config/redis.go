package config

// RedisOptions defines options for building a redis client.
type RedisOptions struct {
	Addrs    []string
	Password string
	Database int
}

// DefaultRedisOptions news a RedisOptions.
func DefaultRedisOptions() *RedisOptions {
	return &RedisOptions{}
}
