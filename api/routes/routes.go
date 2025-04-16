package routes

import (
	"github.com/atilasantos/udacity-golang-crm/api/handlers"
	"github.com/gorilla/mux"
)

const (
	customerBase   = "/customers"
	customerWithID = "/customers/{id}"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc(customerBase, handlers.GetCustomers).Methods("GET")
	router.HandleFunc(customerWithID, handlers.GetCustomer).Methods("GET")
	router.HandleFunc(customerBase, handlers.AddCustomer).Methods("POST")
	router.HandleFunc(customerWithID, handlers.DeleteCustomer).Methods("DELETE")
	router.HandleFunc(customerWithID, handlers.UpdateCustomer).Methods("PUT")
}
