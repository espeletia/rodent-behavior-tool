package runner

import (
	"tusk/internal/config"

	"github.com/nextap-solutions/goNextService"
	"go.uber.org/zap"
)

func StartQueue() error {
	configuration := config.LoadConfig()
	tuskComponents, err := setupService(configuration)
	if err != nil {
		return err
	}
	zap.L().Info("starting tusk queue")
	app := goNextService.NewApplications(tuskComponents.queue)
	return app.Run()
}
