package main

import (
	"cachacariaapi/internal/handlers"
	"cachacariaapi/internal/repositories/user"
	user2 "cachacariaapi/internal/usecases/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
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
	handler := handlers.Handlers{UserHandler: userHandler}

	mux := http.NewServeMux()
	handler.RegisterHandlers(mux)

	log.Print("Server is listening on port 8080")

	http.ListenAndServe(":8080", mux)
}
