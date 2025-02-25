package runner

import (
	"encoding/json"
	"ghiaccio/setup"
	"net/http"
	"valentine/internal/config"
	"valentine/internal/handlers"
	"valentine/internal/middleware"
	"valentine/internal/usecases"

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

	router.Use(middleware.Authentication())

	userUsecase := usecases.NewUserUsecase(http.Client{}, configuration.TuskConfig.URL)
	cageUsecase := usecases.NewCageUsecase(http.Client{}, configuration.TuskConfig.URL)
	videoUsecase := usecases.NewVideoUsecase(http.Client{}, configuration.TuskConfig.URL)

	viewHandler := handlers.NewViewHandler(userUsecase, cageUsecase, videoUsecase)

	commonHandler := handlers.NewCommonHandler()

	router.Handle("/", commonHandler.Handle(viewHandler.Render)).Methods("GET")
	router.Handle("/app", commonHandler.Handle(viewHandler.App)).Methods("GET")
	router.Handle("/cage/{id}", commonHandler.Handle(viewHandler.CageView)).Methods("GET")
	router.Handle("/cage/{id}/message/{message}", commonHandler.Handle(viewHandler.MessageView)).Methods("GET")
	router.Handle("/login", commonHandler.Handle(viewHandler.Login)).Methods("GET")
	router.Handle("/about", commonHandler.Handle(viewHandler.About)).Methods("GET")
	router.Handle("/register", commonHandler.Handle(viewHandler.Register)).Methods("GET")
	router.Handle("/register", commonHandler.Handle(viewHandler.HandleRegisterForm)).Methods("POST")
	router.Handle("/login", commonHandler.Handle(viewHandler.HandleLoginForm)).Methods("POST")

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
