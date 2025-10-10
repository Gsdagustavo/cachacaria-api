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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
	//productRepository := repositories.NewMySQLProductRepository(conn)
	//userRepository := repositories.NewMySQLUserRepository(conn)

	// Use Cases
	authUseCases := usecases.NewAuthUseCases(authRepository, crypt)
	//productUseCases := usecases.NewProductUseCases(productRepository)
	//userUseCases := usecases.NewUserUseCases(userRepository)

	// Modules
	healthModule := modules.NewHealthModule()
	authModule := modules.NewAuthModule(*authUseCases)

	// Assign a router to the server
	cfg.Server.Router = mux.NewRouter()

	cfg.Server.Router.Use()
	mux.CORSMethodMiddleware(r.router)

	// Register routes
	cfg.Server.RegisterModules(healthModule, authModule)

	log.Printf("server running on port %d", cfg.Server.Port)

	err = cfg.Server.Run(*cfg)
	if err != nil {
		panic(err)
	}
}
