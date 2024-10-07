package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceConfig ServiceConfig
	// TracerConfig   TracerConfig
	LoggingConfig LoggingConfig
	// HealthzConfig  HealthzConfig
	// MetricszConfig MetricszConfig
	NatsConfig NatsConfig
}

func LoadConfig(serviceName string) Config {
	config := Config{
		ServiceConfig: loadServiceConfig(serviceName),
		// TracerConfig:   loadTracerConfig(serviceName),
		LoggingConfig: loadLoggingConfig(serviceName),
		// HealthzConfig:  loadHealthzConfig(serviceName),
		// MetricszConfig: loadMetricszConfig(serviceName),
		NatsConfig: LoadNatsConfig(serviceName),
	}

	return config
}

func configViper(configName string, serviceName string) *viper.Viper {
	err := viper.BindEnv(fmt.Sprintf("%s_CONF_PATH", strings.ToUpper(serviceName)))
	if err != nil {
		panic(fmt.Errorf("Fatal error config viper: %w \n", err))
	}
	v := viper.New()
	v.SetEnvPrefix(serviceName)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")

	v.AddConfigPath(fmt.Sprintf("../../%s/configurations/", strings.ToLower(serviceName)))
	v.AddConfigPath(fmt.Sprintf("/app/%s/configurations/", strings.ToLower(serviceName)))
	v.AddConfigPath(viper.GetString(fmt.Sprintf("%s_conf_path", strings.ToLower(serviceName))))
	v.AddConfigPath("./configurations/")
	v.AddConfigPath("/app/configurations/")
	return v
}
