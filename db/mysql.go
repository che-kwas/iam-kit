// Package db is a db instance builder.
package db // import "github.com/che-kwas/iam-kit/db"

import (
	"fmt"

	"github.com/che-kwas/iam-kit/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysqlIns builds a gorm db instance.
func NewMysqlIns(opts *config.MysqlOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
		opts.Username,
		opts.Password,
		opts.Addr,
		opts.Database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.MaxConnLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	return db, nil
}
