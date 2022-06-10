// Package mongo is the mongo instance builder.
package mongo // import "github.com/che-kwas/iam-kit/mongo"

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKey = "mongo"

	defaultMaxPoolSize = 100
)

// MongoOptions defines options for building a mongo client.
type MongoOptions struct {
	URI         string
	MaxPoolSize uint64 `mapstructure:"max-pool-size"`
}

// NewMongoIns creates a mongo client.
func NewMongoIns(ctx context.Context) (*mongo.Client, error) {
	opts, err := getMongoOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new mongo instance with options: %+v", opts)

	mgoOpts := options.Client().ApplyURI(opts.URI).SetMaxPoolSize(opts.MaxPoolSize)

	return mongo.Connect(ctx, mgoOpts)
}

func getMongoOpts() (*MongoOptions, error) {
	opts := &MongoOptions{
		MaxPoolSize: defaultMaxPoolSize,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
