package config

import (
	"errors"
	"strings"

	"github.com/marmotedu/iam/pkg/log"
	"github.com/spf13/viper"
)

var (
	EnvPrefix          = "IAM"
	DefaultConfigPaths = []string{".", "./configs", "/etc/iam"}
)

// LoadConfig loads config from cfgPath or from <DefaultConfigPaths>/<appName>.yaml.
func LoadConfig(cfgPath string, appName string) error {
	log.Infof("Loading config, cfgPath = %s, appName = %s", cfgPath, appName)
	if cfgPath == "" && appName == "" {
		return errors.New("no configuration file specified")
	}

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

	return viper.ReadInConfig()
}
