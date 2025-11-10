package infrastructure

import (
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/config"
	"cachacariaapi/infrastructure/datastore/repositories"
	"cachacariaapi/infrastructure/modules"
	"cachacariaapi/infrastructure/util"
	"log"
	"log/slog"
	"os"

	"github.com/gorilla/mux"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	log.Printf("config file read successfully")

	// Config database
	err = config.Connect(cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Config utils
	maker := util.NewPasetoMaker(cfg.SymmetricKey)
	crypt := util.NewCrypt(maker)

	conn := cfg.Database.Conn

	// Repositories
	authRepository := repositories.NewMySQLAuthRepository(conn)
	userRepository := repositories.NewMySQLUserRepository(conn)
	productRepository := repositories.NewMySQLProductRepository(conn)
	cartRepository := repositories.NewMySQLCartRepository(conn)

	// Use Cases
	authUseCases := usecases.NewAuthUseCases(authRepository, userRepository, crypt)
	userUseCases := usecases.NewUserUseCases(userRepository, authRepository, crypt)
	productUseCases := usecases.NewProductUseCases(productRepository, cfg.Server.BaseURL)
	cartUseCases := usecases.NewCartUseCases(cartRepository, userRepository, productRepository)

	// Modules
	healthModule := modules.NewHealthModule()
	authModule := modules.NewAuthModule(authUseCases)
	userModule := modules.NewUserModule(userUseCases, authUseCases)
	productModule := modules.NewProductModule(productUseCases, crypt)
	cartModule := modules.NewCartModule(cartUseCases, crypt)

	// Assign a router to the server
	router := mux.NewRouter()
	server := &cfg.Server
	server.Router = router

	apiSubrouter := router.PathPrefix("/api").Subrouter()

	// Register health module
	cfg.Server.RegisterModules(server.Router, healthModule)

	// Register modules
	cfg.Server.RegisterModules(apiSubrouter, authModule, userModule, productModule, cartModule)

	log.Printf("server running on port %d", cfg.Server.Port)

	err = cfg.Server.Run(*cfg)

	if err != nil {
		panic(err)
	}
}
