package main

import (
	"log"
	"net/http"

	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/LeandroDeJesus-S/confectionery/internal/routes"
	"github.com/LeandroDeJesus-S/confectionery/internal/config/database"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file:", err)
	}

	db := database.NewDatabaseStarter()
	log.Println("Database started")
	db.MakeMigrations()
	log.Println("All migrations performed")

	baseRouter := mux.NewRouter()
	validator := validator.New(validator.WithRequiredStructEnabled())

	routes.SetupCustomersRoutes(baseRouter, controllers.NewCustomerController(db.DB(), validator))
	log.Println("All routes configured")

	addr := ":8080"
	log.Println("Listening to on address:", addr)
	log.Fatal(http.ListenAndServe(addr, baseRouter))
}
