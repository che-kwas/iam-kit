package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGRPCServer(t *testing.T) {
	assert := assert.New(t)

	server, err := NewGRPCServer()
	assert.Nil(err)
	assert.Equal(DefaultGRPCAddr, server.addr)
}
