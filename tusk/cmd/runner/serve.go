package runner

import (
	"encoding/json"
	"net/http"
	"tusk/internal/config"
	"tusk/internal/handlers"
	"tusk/internal/middleware"
	"tusk/internal/ports/database"
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
	tokenGenerator := tokens.NewTokenGenerator(configuration.JWTConfig.Signature, configuration.JWTConfig.Expiration)
	// placeStore := database.NewDatabasePlaceStore(dbconn)
	userStore := database.NewUserDatabaseStore(dbconn)
	userUsecase := usecases.NewUserUsecase(userStore)
	// placeUsecase := usecases.NewPlaceUsecase(placeStore)
	authUsecase := usecases.NewAuthUsecase(userUsecase, tokenGenerator)
	userHandler := handlers.NewUserHandler(userUsecase, authUsecase)
	// placeHandler := handlers.NewPlaceHandler(placeUsecase)

	router := mux.NewRouter()
	router.Use(middleware.Authentication(authUsecase))
	router.Handle("/", userHandler.Ping()).Methods("GET")
	router.Handle("/me", userHandler.Me()).Methods("GET")
	// router.Handle("/places", placeHandler.GetPlacesByViewport()).Methods("POST")
	router.Handle("/login", userHandler.Login()).Methods("POST")
	// router.Handle("/users", userHandler.GetUsersByUsernamePattern()).Methods("GET")
	router.Handle("/register", userHandler.CreateUser()).Methods("PUT")
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

	return &TuskServiceComponents{
		httpServer: httpComponent,
		cleanup: components.NewLifecycleComponent([]components.LifeCycleFunc{},
			lifecycleRun, nil),
	}, nil

}
