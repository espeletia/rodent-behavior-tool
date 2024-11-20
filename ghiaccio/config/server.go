package config

import (
	"fmt"
	"strings"
	"time"
)

type CorsConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	AllowedHeaders   []string
}

type ServerConfig struct {
	Port            string
	Enabled         bool
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	HandlersTimeout time.Duration

	Cors *CorsConfig
}

func LoadServerConfig(serviceName string) ServerConfig {
	serverConfig := &ServerConfig{}
	v := configViper("server", serviceName)
	err := v.BindEnv("Port", "PORT", fmt.Sprintf("%s_PORT", strings.ToUpper(serviceName)))
	if err != nil {
		panic(err)
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(serverConfig)
	if err != nil {
		panic(err)
	}

	validateServerConfigValues(*serverConfig)

	return *serverConfig
}

func LoadGrpconfig(serviceName string) ServerConfig {
	serverConfig := &ServerConfig{}
	v := configViper("grpc", serviceName)
	err := v.BindEnv("Port", "GRPC_PORT", fmt.Sprintf("%s_GRPC_PORT", strings.ToUpper(serviceName)))
	if err != nil {
		panic(err)
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(serverConfig)
	if err != nil {
		panic(err)
	}

	validateServerConfigValues(*serverConfig)

	return *serverConfig
}

func validateServerConfigValues(serverConfig ServerConfig) ServerConfig {
	if serverConfig.HandlersTimeout > 0 && serverConfig.HandlersTimeout >= serverConfig.WriteTimeout {
		// Otherwise, the server will not be able to respond to the client and response will return EOF
		panic("HandlersTimeout must be less than WriteTimeout")
	}

	return serverConfig
}
