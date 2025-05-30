package main

import (
	"log"

	"github.com/LeandroDeJesus-S/confectionery/config/database"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file:", err)
	}

	db := database.NewDatabaseStarter()
	db.MakeMigrations()
}
