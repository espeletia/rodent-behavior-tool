package config

import "fmt"

type DBConfig struct {
	ConnectionURI string `json:"-"`
	DriverName    string
	RunMigration  bool
}

func loadDbConfig() DBConfig {
	dbConfig := &DBConfig{}
	v := configViper("db")
	err := v.BindEnv("ConnectionURI", "DATABASE_URL")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
	err = v.Unmarshal(dbConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
	return *dbConfig
}
