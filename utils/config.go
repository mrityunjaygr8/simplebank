package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"SB_DB_DRIVER"`
	DBSource      string `mapstructure:"SB_DB_SOURCE"`
	ServerAddress string `mapstructure:"SB_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("sb")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}
