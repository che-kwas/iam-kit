package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitConfig(t *testing.T) {
	assert := assert.New(t)

	// no config file specified
	err := InitConfig("", "")
	assert.NotNil(err)

	// first initialization
	err = InitConfig("./config_test.yaml", "")
	assert.Nil(err)

	cfg := Cfg()
	assert.Equal("127.0.0.1:7777", cfg.HTTPOpts.Addr)
	assert.True(cfg.HTTPOpts.Healthz)
	assert.Equal([]string{"127.0.0.1:6379"}, cfg.RedisOpts.Addrs)
	assert.Equal("", cfg.RedisOpts.Password)

	// ignore second initialization
	err = InitConfig("", "wrongAppName")
	assert.Nil(err)
	assert.Equal("127.0.0.1:7777", Cfg().HTTPOpts.Addr)

}
