package routes

import (
	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/gorilla/mux"
)

func SetupOrdersRoutes(baseRouter *mux.Router, c *controllers.OrdersController) {
	r := baseRouter.PathPrefix("/orders").Subrouter()

	r.HandleFunc("/", c.GetOrders).Methods("GET")
	r.HandleFunc("/", c.CreateOrder).Methods("POST")

	r.HandleFunc("/{id}", c.GetOrder).Methods("GET")
	r.HandleFunc("/{id}", c.UpdateOrder).Methods("PATCH")
	r.HandleFunc("/{id}", c.DeleteOrder).Methods("DELETE")
}
