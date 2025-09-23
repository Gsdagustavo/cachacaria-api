package main

import (
	"cachacariaapi/internal/domain/usecases/product"
	userusecases "cachacariaapi/internal/domain/usecases/user"
	"cachacariaapi/internal/infrastructure/config"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/handlers"
	"cachacariaapi/internal/interfaces/http/handlers/authhandler"
	"cachacariaapi/internal/interfaces/http/handlers/producthandler"
	"cachacariaapi/internal/interfaces/http/handlers/userhandler"
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
	userUseCases := userusecases.NewUserUseCases(userRepository)
	productUseCases := product.NewProductUseCases(productRepository)

	// HANDLERS
	userHandler := userhandler.NewUserHandler(userUseCases)
	authHandler := authhandler.NewAuthHandler(userUseCases)
	productHandler := producthandler.NewProductHandler(productUseCases)

	h := handlers.NewHandlers(userHandler, authHandler, productHandler)
	router := handlers.NewMuxRouter(serverConfig)
	router.StartServer(h, serverConfig)
}
