package runner

import (
	// "context"
	"context"
	"echoes/internal/config"
	"echoes/internal/ports/filemanager"
	"echoes/internal/ports/natsqueue"
	"echoes/internal/setup"
	"echoes/internal/usecases/encoding"
	"echoes/internal/usecases/encoding/video"
	"encoding/json"
	"net/http"

	// "echoes/internal/usecases/encoding/video"
	// "fmt"

	commonSetup "ghiaccio/setup"

	"github.com/nextap-solutions/goNextService"
	"github.com/nextap-solutions/goNextService/components"
)

type EchoesServerComponents struct {
	queue goNextService.Component
}

func Serve() error {
	configuration := config.LoadConfig()
	components, err := setupService(configuration)
	if err != nil {
		return err
	}
	app := goNextService.NewApplications(components.queue)
	return app.Run()
}

func setupService(configuration *config.Config) (*EchoesServerComponents, error) {
	logger := commonSetup.InitLogger(configuration.CommonConfig)
	s, _ := json.MarshalIndent(configuration, "", "\t")
	logger.Info(string(s))

	queue, err := natsqueue.NewNatsQueue(configuration.NatsConfig)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	s3client, err := setup.SetupS3Client(configuration.S3Config, httpClient)
	if err != nil {
		return nil, err
	}

	fileManager := filemanager.NewS3FileManager(s3client)
	mp4Encoder := video.NewVideoMediaEncoder(configuration.EncodingConfig.FfmpegPath, configuration.EncodingConfig.FfprobePath)
	worker := encoding.NewQueueConsumer(fileManager, mp4Encoder, "", configuration.S3Config.URL, configuration.S3Config.Bucket, queue)
	queueComponent := components.NewQueueComponent([]components.QueueHandler{
		func(c chan error) error {
			return queue.HandleVideoJob(context.Background(), worker.ProcessVideoQueue, c)
		},
	},
		components.WithQueueClose(func(ctx context.Context) error {
			err := queue.Close(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	)
	return &EchoesServerComponents{
		queue: queueComponent,
	}, nil
}
