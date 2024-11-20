package config

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	Signature  string
	Expiration time.Duration
}

func loadJWTConfig() JWTConfig {
	JWTConfig := &JWTConfig{}
	v := configViper("jwt")
	err := v.BindEnv("JWT_SIGNATURE", "JWT_EXPIRATION")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(JWTConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return *JWTConfig
}
