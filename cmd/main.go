package main

import (
	"log"
	"net/http"

	"github.com/atilasantos/udacity-golang-crm/api/routes"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
