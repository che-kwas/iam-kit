package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_NoFileSpecified(t *testing.T) {
	assert := assert.New(t)

	err := LoadConfig("", "")
	assert.NotNil(err)
}

func TestLoadConfig_FromCfgPath(t *testing.T) {
	assert := assert.New(t)

	err := LoadConfig("./config_test.yaml", "")
	assert.Nil(err)
	assert.Equal("127.0.0.1:7777", viper.GetString("server.address"))
}

func TestLoadConfig_FromCfgName(t *testing.T) {
	assert := assert.New(t)

	err := LoadConfig("", "config_test.yaml")
	assert.Nil(err)
	assert.Equal("127.0.0.1:7777", viper.GetString("server.address"))
}
