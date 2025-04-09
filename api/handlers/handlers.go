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

func (customer Customer) findById() bool {
	var exists bool = false
	for _, c := range customers {
		if c.ID == customer.ID {
			exists = true
		}
	}

	return exists == true
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

	var customer Customer
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, &customer)
	if customer.findById() {
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
	} else {
		customers = append(customers, customer)
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

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, CONTENT_VALUE)
	var exists bool = true

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for index, customer := range customers {
		if customer.ID == id {
			exists = true
			customers = append(customers[:index], customers[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(customer)
			break
		} else {
			exists = false
		}
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		response, err := json.Marshal(struct {
			Status string `json:"status"`
		}{
			Status: fmt.Sprintf("Customer with id %d does not exists", id),
		})
		if err != nil {
			log.Fatal(err)
		}

		var message map[string]string
		json.Unmarshal(response, &message)
		json.NewEncoder(w).Encode(message)

	}
}
