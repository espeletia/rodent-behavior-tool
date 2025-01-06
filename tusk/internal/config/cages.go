package config

import (
	"fmt"
)

type CagesConfig struct {
	ActivationCodeLength int64
	SecretTokenLength    int64
}

func loadCagesConfig() CagesConfig {
	CagesConfig := &CagesConfig{}
	v := configViper("cages")
	err := v.BindEnv("CAGE_ACTIVATION_CODE_LENGTH")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(CagesConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
	return *CagesConfig
}
