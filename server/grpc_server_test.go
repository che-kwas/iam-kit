package server

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewGRPCServerBuilder_NoConfig(t *testing.T) {
	assert := assert.New(t)

	b := NewGRPCServerBuilder()
	assert.Equal(DefaultGRPCAddr, b.Address)
	assert.Equal(DefaultGRPCMaxMsgSize, b.MaxMsgSize)
	assert.Nil(b.err)
}

func TestNewGRPCServerBuilder_HasConfig(t *testing.T) {
	assert := assert.New(t)

	addr := "127.0.0.1:8888"
	viper.Reset()
	viper.Set("grpc.address", addr)

	b := NewGRPCServerBuilder()
	assert.Equal(addr, b.Address)
	assert.Equal(DefaultGRPCMaxMsgSize, b.MaxMsgSize)
	assert.Nil(b.err)
}

func TestBuild_GRPCServerBuilderError(t *testing.T) {
	assert := assert.New(t)
	b := NewGRPCServerBuilder()
	b.err = errors.New("some error")

	server, err := b.Build()
	assert.Nil(server)
	assert.NotNil(err)
}

func TestBuild_GRPCServerBuilderOk(t *testing.T) {
	assert := assert.New(t)
	viper.Reset()

	server, err := NewGRPCServerBuilder().Build()
	assert.Nil(err)
	assert.Equal(DefaultGRPCAddr, server.address)
}
