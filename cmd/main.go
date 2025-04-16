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

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
		log.Println("APP_PORT not set, defaulting to :3000")
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}
