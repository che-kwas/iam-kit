package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewHTTPServer(t *testing.T) {
	assert := assert.New(t)

	server, err := NewHTTPServer()
	assert.Nil(err)
	assert.Equal(defaultHTTPAddr, server.addr)
	assert.Equal([]string{}, server.middlewares)
	assert.Equal(defaultHTTPPingTimeout, server.pingTimeout)
	assert.Equal(defaultHealthz, server.healthz)
	assert.Equal(defaultMetrics, server.metrics)
	assert.Equal(defaultProfiling, server.profiling)
}
