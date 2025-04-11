package routes

import (
	"fmt"

	"github.com/atilasantos/udacity-golang-crm/api/handlers"
	"github.com/gorilla/mux"
)

const customerRoute = "/customers"

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc(customerRoute, handlers.GetCustomers).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/{id}", customerRoute), handlers.GetCustomer).Methods("GET")
	router.HandleFunc(customerRoute, handlers.AddCustomer).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", customerRoute), handlers.DeleteCustomer).Methods("DELETE")
	router.HandleFunc(customerRoute, handlers.UpdateCustomer).Methods("PUT")
}
