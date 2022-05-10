package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ConfigDBKey = "mysql"

	DefaultHost                  = "127.0.0.1:3306"
	DefaultMaxIdleConnections    = 100
	DefaultMaxOpenConnections    = 100
	DefaultMaxConnectionLifeTime = "10s"
)

// DBBuilder defines options for building a db instance.
type DBBuilder struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int    `mapstructure:"max-idle-connections"`
	MaxOpenConnections    int    `mapstructure:"max-open-connections"`
	MaxConnectionLifeTime string `mapstructure:"max-connection-life-time"`

	err error
}

// NewDBBuilder is used to build a db instance.
func NewDBBuilder() *DBBuilder {
	b := &DBBuilder{
		Host:                  DefaultHost,
		MaxIdleConnections:    DefaultMaxIdleConnections,
		MaxOpenConnections:    DefaultMaxOpenConnections,
		MaxConnectionLifeTime: DefaultMaxConnectionLifeTime,
	}

	b.err = viper.UnmarshalKey(ConfigDBKey, b)
	return b
}

// Build builds a gorm db instance.
func (b *DBBuilder) Build() (*gorm.DB, error) {
	if b.err != nil {
		return nil, b.err
	}

	var maxConnLifeTime time.Duration
	maxConnLifeTime, b.err = time.ParseDuration(b.MaxConnectionLifeTime)
	if b.err != nil {
		return nil, b.err
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
		b.Username,
		b.Password,
		b.Host,
		b.Database)

	var db *gorm.DB
	db, b.err = gorm.Open(mysql.Open(dsn))
	if b.err != nil {
		return nil, b.err
	}

	var sqlDB *sql.DB
	sqlDB, b.err = db.DB()
	if b.err != nil {
		return nil, b.err
	}

	sqlDB.SetMaxOpenConns(b.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	sqlDB.SetMaxIdleConns(b.MaxIdleConnections)

	return db, nil
}
