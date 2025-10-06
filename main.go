package main

import (
	usecases2 "cachacariaapi/domain/usecases"
	config2 "cachacariaapi/infrastructure/config"
	"cachacariaapi/infrastructure/datastore/repositories"
	handlers2 "cachacariaapi/interfaces/http/handlers"
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

	dbConfig := config2.NewDBConfig()

	serverConfig := config2.NewServerConfig(dbConfig)

	dsn := dbConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// REPOSITORIES
	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewMySQLProductRepository(db)

	// USECASES
	userUseCases := usecases2.NewUserUseCases(userRepository)
	productUseCases := usecases2.NewProductUseCases(productRepository)

	// HANDLERS
	userHandler := handlers2.NewUserHandler(userUseCases)
	authHandler := handlers2.NewAuthHandler(userUseCases)
	productHandler := handlers2.NewProductHandler(productUseCases)

	h := handlers2.NewHandlers(userHandler, authHandler, productHandler)
	router := handlers2.NewMuxRouter(serverConfig)
	router.StartServer(h, serverConfig)
}
