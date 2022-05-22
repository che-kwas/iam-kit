package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfig_NoFileSpecified(t *testing.T) {
	assert := assert.New(t)

	cfg, err := NewConfig("", "")
	assert.Nil(cfg)
	assert.NotNil(err)

}

func Test_NewConfig_FromCfgPath(t *testing.T) {
	assert := assert.New(t)

	cfg, err := NewConfig("./config_test.yaml", "")
	assert.Nil(err)

	assert.Equal("127.0.0.1:7777", cfg.HTTPOpts.Addr)
	assert.True(cfg.HTTPOpts.Healthz)
	assert.Equal([]string{"127.0.0.1:6379"}, cfg.RedisOpts.Addrs)
	assert.Equal("", cfg.RedisOpts.Password)
}

func Test_NewConfig_FromAppName(t *testing.T) {
	assert := assert.New(t)

	cfg, err := NewConfig("", "config_test")
	assert.Nil(err)

	assert.Equal("127.0.0.1:7777", cfg.HTTPOpts.Addr)
	assert.True(cfg.HTTPOpts.Healthz)
	assert.Equal([]string{"127.0.0.1:6379"}, cfg.RedisOpts.Addrs)
	assert.Equal("", cfg.RedisOpts.Password)
}
