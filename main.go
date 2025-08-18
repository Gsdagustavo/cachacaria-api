package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {

	cfg := mysql.Config{
		User:   "root",
		Passwd: "admin",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "mydb",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Printf("Database connection established.\n")

	res, err := db.Exec("SELECT * FROM Usuarios")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Res: %v", res)
}
