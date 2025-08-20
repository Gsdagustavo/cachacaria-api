package main

import (
	"cachacariaapi/internal/http/handlers"
	"cachacariaapi/internal/http/handlers/userhandler"
	"cachacariaapi/internal/repositories/userrepository"
	"cachacariaapi/internal/usecases/userusecases"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

/*
database variables

TODO: save into env variables in Docker
*/
var user = "root"
var passwd = "root"
var host = "localhost"
var dbPort = "3306"
var net = "tcp"
var addr = fmt.Sprintf("%s:%s", host, dbPort)
var dbName = "cachacadb"

// http server variables
var serverPort = "8080"

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
