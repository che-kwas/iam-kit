package server

import (
	"testing"

	"github.com/che-kwas/iam-kit/config"
	"github.com/stretchr/testify/assert"
)

func Test_NewGRPCServer(t *testing.T) {
	assert := assert.New(t)

	opts := config.DefaultGRPCOptions()
	server := NewGRPCServer(opts)
	assert.Equal(opts.Addr, server.addr)
}
