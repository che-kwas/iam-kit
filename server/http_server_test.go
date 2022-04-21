package server

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPServerBuilder_NoConfig(t *testing.T) {
	assert := assert.New(t)

	b := NewHTTPServerBuilder()
	assert.Equal(DefaultMode, b.Mode)
	assert.Equal(DefaultHTTPAddr, b.Address)
	assert.Equal([]string{}, b.Middlewares)
	assert.Equal(DefaultHealthz, b.Healthz)
	assert.Equal(DefaultMetrics, b.Metrics)
	assert.Equal(DefaultProfiling, b.Profiling)
	assert.Equal(DefaultPingTimeout, b.PingTimeout)
	assert.Nil(b.err)
}

func TestNewHTTPServerBuilder_HasConfig(t *testing.T) {
	assert := assert.New(t)

	addr := "127.0.0.1:7777"
	viper.Reset()
	viper.Set("http.healthz", false)
	viper.Set("http.address", addr)

	b := NewHTTPServerBuilder()
	assert.False(b.Healthz)
	assert.Equal(addr, b.Address)

	assert.Equal(DefaultMode, b.Mode)
	assert.Equal([]string{}, b.Middlewares)
	assert.Equal(DefaultMetrics, b.Metrics)
	assert.Equal(DefaultProfiling, b.Profiling)
	assert.Equal(DefaultPingTimeout, b.PingTimeout)
	assert.Nil(b.err)
}

func TestBuild_HTTPServerBuilderError(t *testing.T) {
	assert := assert.New(t)
	b := NewHTTPServerBuilder()
	b.err = errors.New("some error")

	server, err := b.Build()
	assert.Nil(server)
	assert.NotNil(err)
}

func TestBuild_HTTPServerBuilderOk(t *testing.T) {
	assert := assert.New(t)
	addr := "127.0.0.1:7777"
	viper.Reset()
	viper.Set("http.address", addr)

	server, err := NewHTTPServerBuilder().Build()
	assert.Nil(err)
	assert.Equal(addr, server.address)
	assert.Equal([]string{}, server.middlewares)
	assert.True(server.healthz)
	assert.False(server.metrics)
	assert.False(server.profiling)
	assert.Equal(10*time.Second, server.pingTimeout)
}
