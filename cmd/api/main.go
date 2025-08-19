package main

import (
	"cachacariaapi/internal/handlers"
	"cachacariaapi/internal/repositories/user"
	user2 "cachacariaapi/internal/usecases/user"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

const port = "8080"

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "admin",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "cachacadb",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// REPOSITORIES
	userRepository := user.NewUserRepository(db)

	// USECASES
	userUseCases := user2.NewUserUseCases(userRepository)

	// HANDLERS
	userHandler := handlers.NewUserHandler(userUseCases)
	h := handlers.Handlers{UserHandler: userHandler}

	router := handlers.NewMuxRouter()
	router.StartServer(h, port)
}
