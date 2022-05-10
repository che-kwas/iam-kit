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

	db              *gorm.DB
	maxConnLifeTime time.Duration
	err             error
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
	if b.err != nil {
		return b
	}

	b.maxConnLifeTime, b.err = time.ParseDuration(b.MaxConnectionLifeTime)
	return b
}

// Build builds a gorm db instance.
func (b *DBBuilder) Build() (*gorm.DB, error) {
	b.openDB().setOptions()

	return b.db, b.err
}

func (b *DBBuilder) openDB() *DBBuilder {
	if b.err != nil {
		return b
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
		b.Username,
		b.Password,
		b.Host,
		b.Database)

	b.db, b.err = gorm.Open(mysql.Open(dsn))
	return b
}

func (b *DBBuilder) setOptions() {
	if b.err != nil {
		return
	}

	var sqlDB *sql.DB
	sqlDB, b.err = b.db.DB()
	if b.err != nil {
		return
	}

	sqlDB.SetMaxOpenConns(b.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(b.maxConnLifeTime)
	sqlDB.SetMaxIdleConns(b.MaxIdleConnections)
}
