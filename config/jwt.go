package config

import "time"

const (
	DefaultJWTTimeout    = time.Duration(24 * time.Second)
	DefaultJWTMaxRefresh = time.Duration(24 * time.Second)
)

// JWTOptions defines options for building a GinJWTMiddleware.
type JWTOptions struct {
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration `mapstructure:"max-refresh"`
}

// DefaultJWTOptions news a JWTOptions.
func DefaultJWTOptions() *JWTOptions {
	return &JWTOptions{
		Timeout:    DefaultJWTTimeout,
		MaxRefresh: DefaultJWTMaxRefresh,
	}
}
