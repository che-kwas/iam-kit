// Package config loads configuration from config file or environment variables.
package config // import "github.com/che-kwas/iam-kit/config"

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	envPrefix          = "IAM"
	defaultConfigPaths = []string{".", "./configs", "/etc/iam"}
)

// LoadConfig loads a config from cfgPath or from <DefaultConfigPaths>/<appName>.yaml.
func LoadConfig(cfgPath, appName string) error {
	log.Printf("Loading config, cfgPath = %s, appName = %s", cfgPath, appName)
	if cfgPath == "" && appName == "" {
		return errors.New("no configuration file specified")
	}

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		for _, path := range defaultConfigPaths {
			viper.AddConfigPath(path)
		}

		viper.SetConfigFile(appName + ".yaml")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return viper.ReadInConfig()
}
