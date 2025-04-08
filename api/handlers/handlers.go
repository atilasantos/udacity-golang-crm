package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Customer struct {
	ID        int
	Name      string
	Role      string
	Email     string
	Phone     int
	Contacted bool
}

func loadCustomers(filePath string) []Customer {
	fileContent, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer fileContent.Close()
	byteResult, _ := io.ReadAll(fileContent)

	var customers []Customer
	json.Unmarshal([]byte(byteResult), &customers)

	return customers

}

var customers []Customer = loadCustomers("api/handlers/customers.json")

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for _, customer := range customers {
		if customer.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(customer)
		}
	}
}
