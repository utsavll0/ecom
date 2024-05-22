package main

import (
	"database/sql"
	"log"

	"github.com/utsavll0/ecom/cmd/api"
	"github.com/utsavll0/ecom/config"
	"github.com/utsavll0/ecom/db"
)

func main() {
	database, err := db.NewPGSQLStorage(config.Envs.DSN)

	if err != nil {
		log.Fatal(err)
	}
	initStorage(database)
	server := api.NewApiServer(":8080", database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(database *sql.DB) {
	err := database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
}
