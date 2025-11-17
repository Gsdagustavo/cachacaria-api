package infrastructure

import (
	"cachacariaapi/domain/usecases"
	util2 "cachacariaapi/domain/util"
	"cachacariaapi/infrastructure/config"
	"cachacariaapi/infrastructure/datastore/repositories"
	"cachacariaapi/infrastructure/modules"
	"cachacariaapi/infrastructure/util"
	"fmt"
	"log"
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
	authManager := util.NewAuthManager(cfg.SymmetricKey)

	conn := cfg.Database.Conn

	// Repositories
	authRepository := repositories.NewMySQLAuthRepository(conn)
	userRepository := repositories.NewMySQLUserRepository(conn)
	productRepository := repositories.NewMySQLProductRepository(conn)
	cartRepository := repositories.NewMySQLCartRepository(conn)

	// Use Cases
	authUseCases := usecases.NewAuthUseCases(authRepository, userRepository, authManager)
	userUseCases := usecases.NewUserUseCases(userRepository, authRepository, authManager)
	productUseCases := usecases.NewProductUseCases(productRepository, cfg.Server.BaseURL)
	cartUseCases := usecases.NewCartUseCases(cartRepository, userRepository, productRepository, cfg.Server.BaseURL)

	// Modules
	healthModule := modules.NewHealthModule()
	authModule := modules.NewAuthModule(authUseCases)
	userModule := modules.NewUserModule(userUseCases, authUseCases)
	productModule := modules.NewProductModule(productUseCases, authManager)
	cartModule := modules.NewCartModule(cartUseCases, authManager)

	// Assign a router to the server
	router := mux.NewRouter()
	server := &cfg.Server
	server.Router = router

	apiSubrouter := router.PathPrefix("/api").Subrouter()

	// Register health module
	cfg.Server.RegisterModules(server.Router, healthModule)

	// Register modules
	cfg.Server.RegisterModules(apiSubrouter, authModule, userModule, productModule, cartModule)

	slog.Info(fmt.Sprintf("server running on port %d", cfg.Server.Port))

	log.Printf("STMP host: %s", cfg.Email.SMTPHost)
	log.Printf("STMP port: %s", cfg.Email.SMTPPort)
	log.Printf("From: %s", cfg.Email.From)
	log.Printf("Username: %s", cfg.Email.Username)
	log.Printf("Password: %s", cfg.Email.Password)

	err = util2.SendEmail(cfg.Email, []string{"gugadanielalvez@gmail.com"}, "Welcome!", "Your account has been created!")
	if err != nil {
		slog.Error("failed to send email", "error", err)
		os.Exit(1)
	}

	err = cfg.Server.Run(*cfg)

	if err != nil {
		panic(err)
	}
}
