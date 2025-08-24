package main

import (
	"cachacariaapi/internal/http/handlers"
	"cachacariaapi/internal/http/handlers/userhandler"
	"cachacariaapi/internal/repositories/userrepository"
	"cachacariaapi/internal/usecases/userusecases"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var user = os.Getenv("DB_USER")
var passwd = os.Getenv("DB_PASSWORD")
var host = os.Getenv("DB_HOST")
var dbPort = os.Getenv("DB_PORT")
var dbName = os.Getenv("DB_NAME")
var net = "tcp"
var addr = fmt.Sprintf("%s:%s", host, dbPort)

var serverPort = os.Getenv("SERVER_PORT")

func main() {
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
	h := handlers.Handlers{UserHandler: userHandler}

	router := handlers.NewMuxRouter()
	router.StartServer(h, serverPort)
}
