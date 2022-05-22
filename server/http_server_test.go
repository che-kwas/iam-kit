package server

import (
	"testing"

	"github.com/che-kwas/iam-kit/config"
	"github.com/stretchr/testify/assert"
)

func Test_NewHTTPServer(t *testing.T) {
	assert := assert.New(t)

	opts := config.DefaultHTTPOptions()
	server := NewHTTPServer(opts)
	assert.Equal(opts.Addr, server.addr)
	assert.Equal(opts.Middlewares, server.middlewares)
	assert.Equal(opts.PingTimeout, server.pingTimeout)
	assert.Equal(opts.Healthz, server.healthz)
	assert.Equal(opts.Metrics, server.metrics)
	assert.Equal(opts.Profiling, server.profiling)
}
