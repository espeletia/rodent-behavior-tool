package config

import (
	commonConfig "ghiaccio/config"
	"github.com/spf13/viper"
)

type Config struct {
	CommonConfig   commonConfig.Config
	EncodingConfig EncodingConfig
	NatsConfig     NatsConfig
	S3Config       S3Config
}

func LoadConfig() *Config {
	return &Config{
		CommonConfig:   commonConfig.LoadConfig("ECHOES"),
		EncodingConfig: loadEncodingConfig(),
		NatsConfig:     loadNatsConfig(),
		S3Config:       loadS3Config(),
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
