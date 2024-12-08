package config

import (
	"fmt"
)

type MigrationsConfig struct {
	MigrationPath     string
	RunNatsMigrations bool
}

func loadMigrationsConfig() MigrationsConfig {
	migrationsConfig := &MigrationsConfig{}
	v := configViper("migrations")
	err := v.BindEnv("RunNatsMigrations", "RUN_NATS_MIGRATIONS")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(migrationsConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return *migrationsConfig
}
