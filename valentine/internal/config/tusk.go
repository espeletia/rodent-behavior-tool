package config

import (
	"fmt"
)

type TuskConfig struct {
	URL string
}

func loadTuskConfig() TuskConfig {
	TuskConfig := &TuskConfig{}
	v := configViper("tusk")
	err := v.BindEnv("URL", "TUSK_URL")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(TuskConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
	return *TuskConfig
}
