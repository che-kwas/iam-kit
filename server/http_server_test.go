package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewHTTPServer(t *testing.T) {
	assert := assert.New(t)

	server, err := NewHTTPServer()
	assert.Nil(err)
	assert.Equal(DefaultHTTPAddr, server.addr)
	assert.Equal([]string{}, server.middlewares)
	assert.Equal(DefaultHTTPPingTimeout, server.pingTimeout)
	assert.Equal(DefaultHealthz, server.healthz)
	assert.Equal(DefaultMetrics, server.metrics)
	assert.Equal(DefaultProfiling, server.profiling)
}
