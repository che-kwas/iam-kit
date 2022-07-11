// Package mongo is the mongo instance builder.
package mongo // import "github.com/che-kwas/iam-kit/mongo"

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKey = "mongo"

	defaultURI         = "mongodb://localhost:27017"
	defaultTimeout     = time.Duration(5 * time.Second)
	defaultMaxPoolSize = 100
)

// MongoOptions defines options for building a mongo client.
type MongoOptions struct {
	URI         string
	Database    string
	Username    string
	Password    string
	Timeout     time.Duration
	MaxPoolSize uint64 `mapstructure:"max-pool-size"`
}

// NewMongoIns creates a mongo client.
func NewMongoIns() (*mongo.Client, error) {
	opts, err := getMongoOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new mongo instance with options: %+v", opts)

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cred := options.Credential{
		AuthSource: opts.Database,
		Username:   opts.Username,
		Password:   opts.Password,
	}

	mgoOpts := options.Client().ApplyURI(opts.URI).SetAuth(cred).
		SetMaxPoolSize(opts.MaxPoolSize)
	cli, err := mongo.Connect(ctx, mgoOpts)
	if err == nil {
		err = cli.Ping(ctx, readpref.Primary())
	}

	if err != nil {
		err = fmt.Errorf("failed to build mongo instance: %s", err.Error())
		return nil, err
	}

	return cli, nil
}

func getMongoOpts() (*MongoOptions, error) {
	opts := &MongoOptions{
		URI:         defaultURI,
		Timeout:     defaultTimeout,
		MaxPoolSize: defaultMaxPoolSize,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
