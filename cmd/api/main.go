package main

import (
	"cachacariaapi/internal/domain/usecases/product"
	"cachacariaapi/internal/domain/usecases/user"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/handlers"
	"cachacariaapi/internal/interfaces/http/handlers/authhandler"
	"cachacariaapi/internal/interfaces/http/handlers/producthandler"
	"cachacariaapi/internal/interfaces/http/handlers/userhandler"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var (
	user       = os.Getenv("DB_USER")
	passwd     = os.Getenv("DB_PASSWORD")
	host       = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbName     = os.Getenv("DB_NAME")
	serverPort = os.Getenv("SERVER_PORT")
)

var net = "tcp"
var addr = fmt.Sprintf("%s:%s", host, dbPort)

func main() {
	//loadJwtEnv()

	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   addr,
		DBName: dbName,
	}

	dsn := cfg.FormatDSN()

	log.Printf("dsn: %s", dsn)

	db, err := sql.Open("mysql", cfg.FormatDSN())

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

	router := handlers.NewMuxRouter()
	router.StartServer(h, serverPort)
}
