package runner

import (
	"encoding/json"
	"ghiaccio/setup"
	"net/http"
	"valentine/internal/config"

	"github.com/gorilla/mux"
	"github.com/nextap-solutions/goNextService"
	"github.com/nextap-solutions/goNextService/components"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type ValentineServiceComponents struct {
	httpServer goNextService.Component
}

func Serve() error {
	configuration := config.LoadConfig()
	valentineComponents, err := setupService(configuration)
	if err != nil {
		return err
	}
	app := goNextService.NewApplications(valentineComponents.httpServer)
	return app.Run()
}

func setupService(configuration *config.Config) (*ValentineServiceComponents, error) {
	logger := setup.InitLogger(configuration.CommonConfig)
	s, err := json.MarshalIndent(configuration, "", "\t")
	if err != nil {
		logger.Error("Failed to marshal configuration", zap.Error(err))
	}

	logger.Info("logger initialized successfully")
	logger.Info(string(s))

	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods("GET")

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := corsMiddleware.Handler(router)
	api := &http.Server{
		Addr:         "0.0.0.0:" + configuration.ServerConfig.Port,
		ReadTimeout:  configuration.ServerConfig.ReadTimeout,
		WriteTimeout: configuration.ServerConfig.WriteTimeout,
		Handler:      handler,
	}
	httpServer := components.NewHttpComponent(handler, components.WithHttpServer(api))

	return &ValentineServiceComponents{
		httpServer: httpServer,
	}, nil
}
