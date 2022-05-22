// Package config manages the configuration of IAM platform.
package config // import "github.com/che-kwas/iam-kit/config"

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	EnvPrefix          = "IAM"
	DefaultConfigPaths = []string{".", "./configs", "/etc/iam"}
)

// Config defines the structure of the iam configuration.
type Config struct {
	HTTPOpts  *HTTPOptions  `mapstructure:"http"`
	GRPCOpts  *GRPCOptions  `mapstructure:"grpc"`
	JWTOpts   *JWTOptions   `mapstructure:"jwt"`
	MysqlOpts *MysqlOptions `mapstructure:"mysql"`
	RedisOpts *RedisOptions `mapstructure:"redis"`

	err error
}

// NewConfig builds a config from cfgPath or from <DefaultConfigPaths>/<appName>.yaml.
func NewConfig(cfgPath, appName string) (*Config, error) {
	log.Printf("Building config, cfgPath = %s, appName = %s", cfgPath, appName)
	if cfgPath == "" && appName == "" {
		return nil, errors.New("no configuration file specified")
	}

	cfg := &Config{
		HTTPOpts:  DefaultHTTPOptions(),
		GRPCOpts:  DefaultGRPCOptions(),
		JWTOpts:   DefaultJWTOptions(),
		MysqlOpts: DefaultMysqlOptions(),
		RedisOpts: DefaultRedisOptions(),
	}

	cfg.loadConfig(cfgPath, appName).unmarshal()
	return cfg, cfg.err
}

func (c *Config) loadConfig(cfgPath, appName string) *Config {
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		for _, path := range DefaultConfigPaths {
			viper.AddConfigPath(path)
		}

		viper.SetConfigFile(appName + ".yaml")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	c.err = viper.ReadInConfig()

	return c
}

func (c *Config) unmarshal() *Config {
	if c.err != nil {
		return c
	}

	c.err = viper.Unmarshal(c)
	return c
}
