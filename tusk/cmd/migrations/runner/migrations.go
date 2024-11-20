package runner

import (
	// "encoding/json"
	// commonSetup "ghiaccio/setup"

	"github.com/pressly/goose/v3"
	// "go.uber.org/zap"
	"tusk/internal/config"
	"tusk/internal/setup"
)

func RunMigrations() error {
	err := RunDBMigrations()
	if err != nil {
		return err
	}
	return nil
}

func RunDBMigrations() error {
	configuration := config.LoadConfig()
	dbConn, err := setup.SetupDb(configuration)
	if err != nil {
		return err
	}

	if err := goose.Up(dbConn, configuration.MigrationsConfig.MigrationPath); err != nil {
		return err
	}

	return nil
}

// func RunNatsMigrations() error {
// 	configuration := config.LoadConfig()
// 	logger := commonSetup.InitLogger(configuration.CommonConfig)
// 	logger.Info("Here i would put my migrations if i had any")
// 	if err != nil {}
// 	err := commonSetup.MigrateNatsStreams(configuration.CommonConfig.NatsConfig)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
