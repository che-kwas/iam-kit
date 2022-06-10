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

	defaultURI = "mongodb://localhost:27017"
)

// MongoOptions defines options for building a mongo client.
type MongoOptions struct {
	URI string
}

// NewMongoIns creates a mongo client.
func NewMongoIns(ctx context.Context) (*mongo.Client, error) {
	opts, err := getMongoOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new mongo instance with options: %+v", opts)

	mgoOpts := options.Client().ApplyURI(opts.URI)

	return mongo.Connect(ctx, mgoOpts)
}

func getMongoOpts() (*MongoOptions, error) {
	opts := &MongoOptions{
		URI: defaultURI,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
