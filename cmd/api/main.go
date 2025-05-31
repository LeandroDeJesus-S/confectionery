package main

import (
	"log"
	"net/http"

	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/LeandroDeJesus-S/confectionery/internal/routes"
	"github.com/LeandroDeJesus-S/confectionery/internal/validators"
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
	db.MakeMigrations()

	baseRouter := mux.NewRouter()
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterValidation("password", validators.PasswordValidator)

	routes.SetupCustomersRoutes(baseRouter, controllers.NewCustomerController(db.DB(), validator))

	log.Fatal(http.ListenAndServe(":8080", baseRouter))
}
