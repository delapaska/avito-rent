package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/delapaska/avito-rent/cmd/api"
	"github.com/delapaska/avito-rent/configs"
	"github.com/delapaska/avito-rent/db"

	"github.com/joho/godotenv"
)

// @title Avito-Rent API
// @version 1.0
// @description This is an API for a test job in Avito

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		configs.Envs.Host, configs.Envs.DBPort,
		configs.Envs.DBUser, configs.Envs.DBPassword, configs.Envs.DBName)
	log.Println("DB CONN ", psqlInfo)
	db, err := db.NewPostgresSQLStorage(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)
	log.Println("server started on port:", configs.Envs.Port)
	srv := api.NewAPIServer(db)
	srv.Run()

}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
