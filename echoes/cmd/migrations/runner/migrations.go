package runner

import (
	"echoes/internal/config"
	commonSetup "ghiaccio/setup"
)

func RunMigrations() error {
	err := RunNatsMigrations()
	if err != nil {
		return err
	}
	return nil
}

func RunNatsMigrations() error {
	configuration := config.LoadConfig()
	commonSetup.InitLogger(configuration.CommonConfig)

	err := commonSetup.MigrateNatsStreams(configuration.CommonConfig.NatsConfig)
	if err != nil {
		return err
	}

	return nil
}
