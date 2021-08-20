package config

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// InitConfig reads config into viper
func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/dongibot")
	viper.AddConfigPath("config")

	// mapping config keys to env vars
	viper.SetEnvPrefix("DONGIBOT")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func InitLogrus() {
	lvl, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logrus.WithError(err).Fatal("failed to parse log level")
	}

	logrus.SetLevel(lvl)
}
