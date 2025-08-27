package main

import (
	"cachacariaapi/internal/http/handlers"
	"cachacariaapi/internal/http/handlers/authhandler"
	"cachacariaapi/internal/http/handlers/userhandler"
	"cachacariaapi/internal/repositories/userrepository"
	"cachacariaapi/internal/usecases/userusecases"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/lpernett/godotenv"
)

//var user = os.Getenv("DB_USER")
//var passwd = os.Getenv("DB_PASSWORD")
//var host = os.Getenv("DB_HOST")
//var dbPort = os.Getenv("DB_PORT")
//var dbName = os.Getenv("DB_NAME")
//var net = "tcp"
//var addr = fmt.Sprintf("%s:%s", host, dbPort)
//var serverPort = os.Getenv("PORT")

var user = "root"
var passwd = "admin"
var host = "127.0.0.1"
var dbPort = "3306"
var dbName = "cachacadb"
var net = "tcp"
var addr = fmt.Sprintf("%s:%s", host, dbPort)
var serverPort = "8080"

func main() {
	loadJwtEnv()

	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   addr,
		DBName: dbName,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// REPOSITORIES
	userRepository := userrepository.NewUserRepository(db)

	// USECASES
	userUseCases := userusecases.NewUserUseCases(userRepository)

	// HANDLERS
	userHandler := userhandler.NewUserHandler(userUseCases)
	authHandler := authhandler.NewAuthHandler(userUseCases)

	h := handlers.Handlers{UserHandler: userHandler, AuthHandler: authHandler}

	router := handlers.NewMuxRouter()
	router.StartServer(h, serverPort)
}

func loadJwtEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env file could not be loaded. err: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("jwt secret key was not set")
	}
}
