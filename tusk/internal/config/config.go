package config

import (
	commonConfig "ghiaccio/config"
	"github.com/spf13/viper"
)

type Config struct {
	CommonConfig     commonConfig.Config
	ServerConfig     commonConfig.ServerConfig
	S3Config         S3Config
	DBConfig         DBConfig
	MigrationsConfig MigrationsConfig
	HashConfig       HashConfig
	JWTConfig        JWTConfig
}

func LoadConfig() *Config {
	return &Config{
		CommonConfig:     commonConfig.LoadConfig("TUSK"),
		ServerConfig:     commonConfig.LoadServerConfig("TUSK"),
		S3Config:         loadS3Config(),
		DBConfig:         loadDbConfig(),
		MigrationsConfig: loadMigrationsConfig(),
		HashConfig:       loadHashConfig(),
		JWTConfig:        loadJWTConfig(),
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
