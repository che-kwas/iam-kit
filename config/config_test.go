package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig_NoFileSpecified(t *testing.T) {
	assert := assert.New(t)

	err := LoadConfig("", "")
	assert.NotNil(err)

}

func Test_NewConfig_FromCfgPath(t *testing.T) {
	assert := assert.New(t)

	viper.Reset()
	err := LoadConfig("./config_test.yaml", "")
	assert.Nil(err)
	assert.Equal("localhost:7777", viper.Get("http.addr"))
}

func Test_NewConfig_FromAppName(t *testing.T) {
	assert := assert.New(t)

	viper.Reset()
	err := LoadConfig("", "config_test")
	assert.Nil(err)
	assert.Equal("localhost:7777", viper.Get("http.addr"))
}
