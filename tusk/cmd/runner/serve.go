package runner

import (
	"context"
	"encoding/json"
	"net/http"
	"tusk/internal/config"
	"tusk/internal/handlers"
	"tusk/internal/middleware"
	"tusk/internal/ports/database"
	"tusk/internal/ports/filemanager"
	"tusk/internal/ports/natsqueue"
	"tusk/internal/ports/tokens"
	"tusk/internal/setup"
	"tusk/internal/usecases"

	"github.com/gorilla/mux"
	"github.com/nextap-solutions/goNextService"
	"github.com/nextap-solutions/goNextService/components"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type TuskServiceComponents struct {
	httpServer goNextService.Component
	queue      goNextService.Component
	cleanup    goNextService.Component
}

func Serve() error {
	configuration := config.LoadConfig()
	tuskComponents, err := setupService(configuration)
	if err != nil {
		return err
	}
	app := goNextService.NewApplications(tuskComponents.httpServer)
	return app.Run()
}

func setupService(configuration *config.Config) (*TuskServiceComponents, error) {
	logger := setup.InitLogger(*&configuration.CommonConfig)
	s, err := json.MarshalIndent(configuration, "", "\t")
	if err != nil {
		logger.Error("Failed to marshal configuration", zap.Error(err))
		return nil, err
	}

	logger.Info("Logger initialized successfully")
	logger.Info(string(s))
	dbconn, err := setup.SetupDb(configuration)
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

	// token port
	tokenGenerator := tokens.NewTokenGenerator(configuration.JWTConfig.Signature, configuration.JWTConfig.Expiration)

	// s3 file management port
	fileManager := filemanager.NewS3FileManager(s3client)

	// natsqueue port
	natsqueue, err := natsqueue.NewNatsQueue(configuration.NatsConfig)
	if err != nil {
		return nil, err
	}

	// database ports
	userStore := database.NewUserDatabaseStore(dbconn)
	mediaStore := database.NewMediaDatabaseStore(dbconn)
	videoStore := database.NewVideoDatabaseStore(dbconn)

	// usecases
	userUsecase := usecases.NewUserUsecase(userStore)
	mediaUsecase := usecases.NewMediaUsecase(mediaStore, fileManager, configuration.S3Config.URL, configuration.S3Config.UploadsPathPrefix, configuration.S3Config.Bucket)
	videoUsecase := usecases.NewVideoUsecase(mediaUsecase, videoStore, natsqueue)
	authUsecase := usecases.NewAuthUsecase(userUsecase, tokenGenerator)

	// rest handlers
	userHandler := handlers.NewUserHandler(userUsecase, authUsecase)
	mediaHandler := handlers.NewMediaHandler(mediaUsecase)
	videoHandler := handlers.NewVideoAnalysisHandler(videoUsecase)
	commonHandler := handlers.NewCommonHandler()

	router := mux.NewRouter()
	router.Use(middleware.Authentication(authUsecase))
	// connectivity test
	router.Handle("/", commonHandler.Handle(userHandler.Ping)).Methods("GET")

	// users
	router.Handle("/register", commonHandler.Handle(userHandler.CreateUser)).Methods("PUT")
	router.Handle("/login", commonHandler.Handle(userHandler.Login)).Methods("POST")
	router.Handle("/me", commonHandler.Handle(userHandler.Me)).Methods("GET")

	// media
	router.Handle("/upload", commonHandler.Handle(mediaHandler.Upload)).Methods("PUT")

	// videos
	router.Handle("/video", commonHandler.Handle(videoHandler.CreateVideoAnalysis)).Methods("PUT")
	router.Handle("/video/{id}", commonHandler.Handle(videoHandler.GetVideoAnalysisByID)).Methods("GET")

	specHandler, _, err := handlers.HandleSwaggerFile()
	if err != nil {
		logger.Error("Swagger handler disabled, cause:", zap.Error(err))
	} else {
		router.Handle("/swagger.json", specHandler).Methods("GET")
		router.Handle("/swagger", handlers.HandleSwaggerUI()).Methods("GET")
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(router)
	api := http.Server{
		Addr:         "0.0.0.0:" + configuration.ServerConfig.Port,
		ReadTimeout:  configuration.ServerConfig.ReadTimeout,
		WriteTimeout: configuration.ServerConfig.WriteTimeout,
		Handler:      handler,
	}
	httpComponent := components.NewHttpComponent(handler, components.WithHttpServer(&api))
	var lifecycleRun components.LifeCycleFunc

	queueComponent := components.NewQueueComponent([]components.QueueHandler{
		func(c chan error) error {
			return natsqueue.HandleAnalystJobResult(context.Background(), videoUsecase.ProcessAnalystJobResultQueue, c)
		},
		func(c chan error) error {
			return natsqueue.HandleEncodingJobResult(context.Background(), videoUsecase.ProcessEncodingJobResultQueue, c)
		},
	},
		components.WithQueueClose(func(ctx context.Context) error {
			err := natsqueue.Close(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	)

	return &TuskServiceComponents{
		httpServer: httpComponent,
		cleanup: components.NewLifecycleComponent([]components.LifeCycleFunc{},
			lifecycleRun, nil),
		queue: queueComponent,
	}, nil

}
