package routes

import (
	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/gorilla/mux"
)

func SetupCakeRoutes(baseRouter *mux.Router, cakeController *controllers.CakeController) {
	customersRouter := baseRouter.PathPrefix("/cakes").Subrouter()

	customersRouter.HandleFunc("/", cakeController.GetCakes).Methods("GET")
	customersRouter.HandleFunc("/", cakeController.CreateCake).Methods("POST")

	customersRouter.HandleFunc("/{id}", cakeController.GetCake).Methods("GET")
	customersRouter.HandleFunc("/{id}", cakeController.UpdateCake).Methods("PATCH")
	customersRouter.HandleFunc("/{id}", cakeController.DeleteCake).Methods("DELETE")
}
