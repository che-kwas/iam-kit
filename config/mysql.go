package config

import (
	"time"
)

const (
	DefaultMysqlAddr            = "127.0.0.1:3306"
	DefaultMysqlMaxIdleConns    = 100
	DefaultMysqlMaxOpenConns    = 100
	DefaultMysqlMaxConnLifeTime = time.Duration(10 * time.Second)
)

// MysqlOptions defines options for building a mysql instance.
type MysqlOptions struct {
	Addr            string
	Username        string
	Password        string
	Database        string
	MaxIdleConns    int           `mapstructure:"max-idle-conns"`
	MaxOpenConns    int           `mapstructure:"max-open-conns"`
	MaxConnLifeTime time.Duration `mapstructure:"max-connection-life-time"`
}

// DefaultMysqlOptions news an MySQLOptions.
func DefaultMysqlOptions() *MysqlOptions {
	return &MysqlOptions{
		Addr:            DefaultMysqlAddr,
		MaxIdleConns:    DefaultMysqlMaxIdleConns,
		MaxOpenConns:    DefaultMysqlMaxOpenConns,
		MaxConnLifeTime: DefaultMysqlMaxConnLifeTime,
	}
}
