package routes

import (
	"github.com/atilasantos/udacity-golang-crm/api/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/customers", handlers.GetCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", handlers.GetCustomer).Methods("GET")
	router.HandleFunc("/customers", handlers.AddCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", handlers.DeleteCustomer).Methods("DELETE")
}
