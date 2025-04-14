package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atilasantos/udacity-golang-crm/api/routes"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router); err != nil {
		log.Fatal(err)
	}
}
