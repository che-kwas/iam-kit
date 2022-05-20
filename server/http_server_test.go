package server

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewHTTPServerBuilder_NoConfig(t *testing.T) {
	assert := assert.New(t)

	b := NewHTTPServerBuilder()
	assert.Equal(DefaultMode, b.Mode)
	assert.Equal(DefaultHTTPAddr, b.Addr)
	assert.Equal([]string{}, b.Middlewares)
	assert.Equal(DefaultHealthz, b.Healthz)
	assert.Equal(DefaultMetrics, b.Metrics)
	assert.Equal(DefaultProfiling, b.Profiling)
	assert.Equal(DefaultPingTimeout, b.PingTimeout)
	assert.Nil(b.err)
}

func Test_NewHTTPServerBuilder_HasConfig(t *testing.T) {
	assert := assert.New(t)

	addr := "127.0.0.1:7777"
	viper.Reset()
	viper.Set("http.healthz", false)
	viper.Set("http.addr", addr)

	b := NewHTTPServerBuilder()
	assert.False(b.Healthz)
	assert.Equal(addr, b.Addr)

	assert.Equal(DefaultMode, b.Mode)
	assert.Equal([]string{}, b.Middlewares)
	assert.Equal(DefaultMetrics, b.Metrics)
	assert.Equal(DefaultProfiling, b.Profiling)
	assert.Equal(DefaultPingTimeout, b.PingTimeout)
	assert.Nil(b.err)
}

func Test_Build_HTTPServerBuilderError(t *testing.T) {
	assert := assert.New(t)
	b := NewHTTPServerBuilder()
	b.err = errors.New("some error")

	server, err := b.Build()
	assert.Nil(server)
	assert.NotNil(err)
}

func Test_Build_HTTPServerBuilderOk(t *testing.T) {
	assert := assert.New(t)
	addr := "127.0.0.1:7777"
	viper.Reset()
	viper.Set("http.addr", addr)

	server, err := NewHTTPServerBuilder().Build()
	assert.Nil(err)
	assert.Equal(addr, server.addr)
	assert.Equal([]string{}, server.middlewares)
	assert.True(server.healthz)
	assert.False(server.metrics)
	assert.False(server.profiling)
	assert.Equal(10*time.Second, server.pingTimeout)
}
