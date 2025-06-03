package main

import (
	"log"
	"net/http"

	"github.com/LeandroDeJesus-S/confectionery/internal/config/database"
	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/LeandroDeJesus-S/confectionery/internal/routes"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load .env file:", err)
	}

	dbStarter := database.NewDatabaseStarter()
	db := dbStarter.DB()
	log.Println("Database started")

	dbStarter.MakeMigrations()
	log.Println("All migrations performed")

	baseRouter := mux.NewRouter()
	validator := validator.New(validator.WithRequiredStructEnabled())

	routes.SetupCustomersRoutes(baseRouter, controllers.NewCustomerController(db, validator))
	routes.SetupCakeRoutes(baseRouter, controllers.NewCakeController(db, validator))
	routes.SetupOrdersRoutes(baseRouter, controllers.NewOrdersController(db, validator))
	log.Println("All routes configured")

	addr := ":8080"
	log.Println("Listening to on address:", addr)
	log.Fatal(http.ListenAndServe(addr, baseRouter))
}
