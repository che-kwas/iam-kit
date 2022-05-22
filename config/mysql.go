package config

import (
	"time"
)

const (
	DefaultMysqlAddr            = "127.0.0.1:3306"
	DefaultMysqlMaxIdleConn     = 100
	DefaultMysqlMaxOpenConn     = 100
	DefaultMysqlMaxConnLifeTime = time.Duration(10 * time.Second)
)

// MysqlOptions defines options for building a mysql instance.
type MysqlOptions struct {
	Addr                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int           `mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `mapstructure:"max-connection-life-time"`
}

// DefaultMysqlOptions news an MySQLOptions.
func DefaultMysqlOptions() *MysqlOptions {
	return &MysqlOptions{
		Addr:                  DefaultMysqlAddr,
		MaxIdleConnections:    DefaultMysqlMaxIdleConn,
		MaxOpenConnections:    DefaultMysqlMaxOpenConn,
		MaxConnectionLifeTime: DefaultMysqlMaxConnLifeTime,
	}
}
