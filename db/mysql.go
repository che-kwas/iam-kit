package db

// import (
// 	"fmt"
// 	"time"

// 	"github.com/spf13/viper"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// )

// const (
// 	ConfigDBKey = "mysql"

// 	DefaultHost                  = "127.0.0.1:3306"
// 	DefaultMaxIdleConnections    = 100
// 	DefaultMaxOpenConnections    = 100
// 	DefaultMaxConnectionLifeTime = "10s"
// 	DefaultSlowThreshold         = "1s"

// 	// DefaultLogLevel defines default GORM log level
// 	// 1: silent, 2:error, 3:warn, 4:info
// 	DefaultLogLevel = 2
// )

// // DBBuilder defines options for building a db instance.
// type DBBuilder struct {
// 	Host                  string
// 	Username              string
// 	Password              string
// 	Database              string
// 	MaxIdleConnections    int    `mapstructure:"max-idle-connections"`
// 	MaxOpenConnections    int    `mapstructure:"max-open-connections"`
// 	MaxConnectionLifeTime string `mapstructure:"max-connection-life-time"`
// 	SlowThreshold         string `mapstructure:"slow-threshold"`
// 	LogLevel              int    `mapstructure:"log-level"`
// 	Logger                logger.Interface

// 	err error
// }

// // NewDBBuilder is used to build a db instance.
// func NewDBBuilder(myLogger logger.Interface) *DBBuilder {
// 	b := &DBBuilder{
// 		Host:                  DefaultHost,
// 		MaxIdleConnections:    DefaultMaxIdleConnections,
// 		MaxOpenConnections:    DefaultMaxOpenConnections,
// 		MaxConnectionLifeTime: DefaultMaxConnectionLifeTime,
// 		SlowThreshold:         DefaultSlowThreshold,
// 		LogLevel:              DefaultLogLevel,
// 		Logger:                myLogger,
// 	}

// 	b.err = viper.UnmarshalKey(ConfigDBKey, b)
// 	return b
// }

// // Build builds a gorm db instance.
// func (b *DBBuilder) Build() (*gorm.DB, error) {
// 	if b.err != nil {
// 		return nil, b.err
// 	}

// 	var maxConnLifeTime time.Duration
// 	maxConnLifeTime, b.err = time.ParseDuration(b.MaxConnectionLifeTime)
// 	if b.err != nil {
// 		return nil, b.err
// 	}

// 	var SlowThreshold time.Duration
// 	SlowThreshold, b.err = time.ParseDuration(b.MaxConnectionLifeTime)
// 	if b.err != nil {
// 		return nil, b.err
// 	}

// 	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
// 		b.Username,
// 		b.Password,
// 		b.Host,
// 		b.Database)

// 	var db *gorm.DB
// 	db, b.err = gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: newLogger(b.Logger, b.SlowThreshold, b.LogLevel),
// 	})
// }

// // New create a new gorm db instance with the given options.
// func New(opts *Options) (*gorm.DB, error) {
// 	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local`,
// 		opts.Username,
// 		opts.Password,
// 		opts.Host,
// 		opts.Database)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: opts.Logger,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
// 	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
// 	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

// 	return db, nil
// }

// func newLogger(myLogger logger.Interface, slowThreshold string, logLevel int) logger.Interface {

// 	return logger.New(
// 		myLogger,
// 		logger.Config{
// 			SlowThreshold:             time.Second,   // 慢 SQL 阈值
// 			LogLevel:                  logger.Silent, // 日志级别
// 			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
// 			Colorful:                  false,         // 禁用彩色打印
// 		},
// 	)
// }
