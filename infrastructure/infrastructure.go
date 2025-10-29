package infrastructure

import (
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/config"
	"cachacariaapi/infrastructure/datastore/repositories"
	"cachacariaapi/infrastructure/modules"
	"cachacariaapi/infrastructure/util"
	"fmt"
	"log/slog"
	"os"

	"github.com/gorilla/mux"
)

func Init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
	}

	slog.Info("config loaded")

	// Config database
	err = config.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
	}

	// Config utils
	maker := util.NewPasetoMaker(cfg.SymmetricKey)
	crypt := util.NewCrypt(maker)

	conn := cfg.Database.Conn

	// Repositories
	authRepository := repositories.NewMySQLAuthRepository(conn)
	userRepository := repositories.NewMySQLUserRepository(conn)
	productRepository := repositories.NewMySQLProductRepository(conn)

	// Use Cases
	authUseCases := usecases.NewAuthUseCases(authRepository, crypt)
	userUseCases := usecases.NewUserUseCases(userRepository)
	productUseCases := usecases.NewProductUseCases(productRepository, cfg.Server.BaseURL)

	// Modules
	healthModule := modules.NewHealthModule()
	authModule := modules.NewAuthModule(authUseCases)
	userModule := modules.NewUserModule(userUseCases)
	productModule := modules.NewProductModule(productUseCases)

	// Assign a router to the server
	router := mux.NewRouter()
	server := &cfg.Server
	server.Router = router

	apiSubrouter := router.PathPrefix("/api").Subrouter()

	// Register health module
	cfg.Server.RegisterModules(server.Router, healthModule)

	// Register modules
	cfg.Server.RegisterModules(apiSubrouter, authModule, userModule, productModule)

	slog.Info(fmt.Sprintf("server running on port %d", cfg.Server.Port))

	err = cfg.Server.Run(*cfg)
	if err != nil {
		panic(err)
	}
}
