package main

import (
	"cachacariaapi/internal/domain/usecases"
	"cachacariaapi/internal/infrastructure/config"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/handlers"
	"database/sql"
	"log"
	"log/slog"
	"os"
)

func main() {
	filePath := "logs/app.log"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	dbConfig := config.NewDBConfig()

	serverConfig := config.NewServerConfig(dbConfig)

	dsn := dbConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// REPOSITORIES
	userRepository := persistence.NewUserRepository(db)
	productRepository := persistence.NewMySQLProductRepository(db)

	// USECASES
	userUseCases := usecases.NewUserUseCases(userRepository)
	productUseCases := usecases.NewProductUseCases(productRepository)

	// HANDLERS
	userHandler := handlers.NewUserHandler(userUseCases)
	authHandler := handlers.NewAuthHandler(userUseCases)
	productHandler := handlers.NewProductHandler(productUseCases)

	h := handlers.NewHandlers(userHandler, authHandler, productHandler)
	router := handlers.NewMuxRouter(serverConfig)
	router.StartServer(h, serverConfig)
}
