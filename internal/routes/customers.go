package routes

import (
	"github.com/LeandroDeJesus-S/confectionery/internal/controllers"
	"github.com/gorilla/mux"
)

func SetupCustomersRoutes(baseRouter *mux.Router, customerController *controllers.CustomerController) {
	customersRouter := baseRouter.PathPrefix("/customers").Subrouter()

	customersRouter.HandleFunc("/", customerController.GetAllCustomers).Methods("GET")
	customersRouter.HandleFunc("/", customerController.CreateCustomer).Methods("POST")

	customersRouter.HandleFunc("/{id}", customerController.GetCustomer).Methods("GET")
	customersRouter.HandleFunc("/{id}", customerController.UpdateCustomer).Methods("PATCH")
	customersRouter.HandleFunc("/{id}", customerController.DeleteCustomer).Methods("DELETE")
}
