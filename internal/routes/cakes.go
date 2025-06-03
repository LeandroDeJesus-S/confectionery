package routes

import (
	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/gorilla/mux"
)

func SetupCakeRoutes(baseRouter *mux.Router, cakeController *controllers.CakeController) {
	r := baseRouter.PathPrefix("/cakes").Subrouter()

	r.HandleFunc("/", cakeController.GetCakes).Methods("GET")
	r.HandleFunc("/", cakeController.CreateCake).Methods("POST")

	r.HandleFunc("/{id}", cakeController.GetCake).Methods("GET")
	r.HandleFunc("/{id}", cakeController.UpdateCake).Methods("PATCH")
	r.HandleFunc("/{id}", cakeController.DeleteCake).Methods("DELETE")
}
