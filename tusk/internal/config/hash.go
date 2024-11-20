package config

import "fmt"

type HashConfig struct {
	Salt string
}

func loadHashConfig() HashConfig {
	HashConfig := &HashConfig{}
	v := configViper("hash")
	err := v.BindEnv("HASH_SALT")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(HashConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
	return *HashConfig
}
