// Package mysql is the mysql instance builder.
package mysql // import "github.com/che-kwas/iam-kit/mysql"

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKey = "mysql"

	defaultAddr            = "127.0.0.1:3306"
	defaultMaxIdleConns    = 100
	defaultMaxOpenConns    = 100
	defaultMaxConnLifeTime = time.Duration(10 * time.Second)
	defaultSlowThreshold   = time.Duration(200 * time.Millisecond)
	// GORM log level, 1: silent, 2:error, 3:warn, 4:info
	defaultLogLevel = 1
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
	SlowThreshold   time.Duration `mapstructure:"slow-threshold"`
	LogLevel        int           `mapstructure:"log-level"`
}

// NewMysqlIns creates a gorm db instance.
func NewMysqlIns() (*gorm.DB, error) {
	opts, err := getMysqlOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new mysql instance with options: %+v", opts)

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
		opts.Username,
		opts.Password,
		opts.Addr,
		opts.Database)

	cfg := &gorm.Config{Logger: newLogger(opts.SlowThreshold, opts.LogLevel)}
	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		err = fmt.Errorf("failed to build mysql instance: %s", err.Error())
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
		Addr:            defaultAddr,
		MaxIdleConns:    defaultMaxIdleConns,
		MaxOpenConns:    defaultMaxOpenConns,
		MaxConnLifeTime: defaultMaxConnLifeTime,
		SlowThreshold:   defaultSlowThreshold,
		LogLevel:        defaultLogLevel,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
