package config

import (
	commonConfig "ghiaccio/config"
	"github.com/spf13/viper"
)

type Config struct {
	CommonConfig commonConfig.Config
	ServerConfig commonConfig.ServerConfig
	TuskConfig   TuskConfig
}

func LoadConfig() *Config {
	return &Config{
		CommonConfig: commonConfig.LoadConfig("VALENTINE"),
		ServerConfig: commonConfig.LoadServerConfig("VALENTINE"),
		TuskConfig:   loadTuskConfig(),
	}
}

func configViper(configName string) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configurations/")
	return v
}
