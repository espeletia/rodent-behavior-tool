package runner

import (
	commonSetup "ghiaccio/setup"

	"github.com/pressly/goose/v3"
	"tusk/internal/config"
	"tusk/internal/setup"
)

func RunMigrations() error {
	configuration := config.LoadConfig()
	err := RunDBMigrations(configuration)
	if err != nil {
		return err
	}
	err = RunNatsMigrations(configuration)
	if err != nil {
		return err
	}
	return nil
}

func RunDBMigrations(configuration *config.Config) error {
	dbConn, err := setup.SetupDb(configuration)
	if err != nil {
		return err
	}

	if err := goose.Up(dbConn, configuration.MigrationsConfig.MigrationPath); err != nil {
		return err
	}

	return nil
}

func RunNatsMigrations(configuration *config.Config) error {
	err := commonSetup.MigrateNatsStreams(configuration.CommonConfig.NatsConfig)
	if err != nil {
		return err
	}

	return nil
}
