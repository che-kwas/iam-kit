// Package mysql defines the global mysql instance.
package mysql // import "github.com/che-kwas/iam-kit/mysql"

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ConfKey = "mysql"

	DefaultAddr            = "127.0.0.1:3306"
	DefaultMaxIdleConns    = 100
	DefaultMaxOpenConns    = 100
	DefaultMaxConnLifeTime = time.Duration(10 * time.Second)
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

// NewMysqlIns creates a gorm db instance.
func NewMysqlIns() (*gorm.DB, error) {
	opts, err := getMysqlOpts()
	if err != nil {
		return nil, err
	}
	log.Printf("NewMysqlIns, opts: %+v", opts)

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
		opts.Username,
		opts.Password,
		opts.Addr,
		opts.Database)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	// set options
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.MaxConnLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)

	return db, nil
}

func getMysqlOpts() (*MysqlOptions, error) {
	opts := &MysqlOptions{
		Addr:            DefaultAddr,
		MaxIdleConns:    DefaultMaxIdleConns,
		MaxOpenConns:    DefaultMaxOpenConns,
		MaxConnLifeTime: DefaultMaxConnLifeTime,
	}

	if err := viper.UnmarshalKey(ConfKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
