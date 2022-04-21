package util

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	EnvPrefix          = "IAM"
	DefaultConfigPaths = []string{".", "./configs", "/etc/iam"}
)

// LoadConfig loads config from file and ENV variables if set.
func LoadConfig(cfgPath string, cfgName string) error {
	log.Debugf("Loading config, cfgPath = %s, cfgName = %s", cfgPath, cfgName)
	if cfgPath == "" && cfgName == "" {
		return errors.New("no configuration file specified")
	}

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		for _, path := range DefaultConfigPaths {
			viper.AddConfigPath(path)
		}

		viper.SetConfigFile(cfgName)
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return viper.ReadInConfig()
}
