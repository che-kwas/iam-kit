// Package config manages the configuration of IAM platform.
package config // import "github.com/che-kwas/iam-kit/config"

import (
	"errors"
	"log"
	"strings"
	"sync"

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

var (
	cfg  *Config
	once sync.Once
)

// Cfg returns the global cfg instance.
func Cfg() *Config {
	if cfg == nil {
		log.Fatal("Configuration not initialized.")
	}

	return cfg
}

// InitConfig initializes the config from cfgPath or from <DefaultConfigPaths>/<appName>.yaml.
func InitConfig(cfgPath, appName string) error {
	log.Printf("Initializing config, cfgPath = %s, appName = %s", cfgPath, appName)
	if cfgPath == "" && appName == "" {
		return errors.New("no configuration file specified")
	}

	if cfg != nil {
		return nil
	}

	var err error
	once.Do(func() {
		cfg = newConfig().loadConfig(cfgPath, appName).unmarshal()
		if err = cfg.err; err != nil {
			cfg = nil
		}
	})

	return err
}

func newConfig() *Config {
	return &Config{
		HTTPOpts:  DefaultHTTPOptions(),
		GRPCOpts:  DefaultGRPCOptions(),
		JWTOpts:   DefaultJWTOptions(),
		MysqlOpts: DefaultMysqlOptions(),
		RedisOpts: DefaultRedisOptions(),
	}
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
