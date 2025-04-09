package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var CONTENT_TYPE string = "Content-Type"
var CONTENT_VALUE string = "application/json"

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

func AddCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, CONTENT_VALUE)
	var customerExists bool = true

	var customer Customer
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, &customer)
	for _, c := range customers {
		if c.ID == customer.ID {
			w.WriteHeader(http.StatusConflict)
			response, err := json.Marshal(struct {
				Status string `json:"status"`
			}{
				Status: fmt.Sprintf("Customer %s already exists", customer.Name),
			})
			if err != nil {
				log.Fatal(err)
			}

			var message map[string]string
			json.Unmarshal(response, &message)
			json.NewEncoder(w).Encode(message)
		}
	}

	if !customerExists {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customer)
	}

}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, CONTENT_VALUE)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, CONTENT_VALUE)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for _, customer := range customers {
		if customer.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(customer)
		}
	}
}
