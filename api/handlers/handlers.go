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

func (customer Customer) findById() (bool, int) {
	var exists bool = false
	var index int
	for i, c := range customers {
		if c.ID == customer.ID {
			index = i
			exists = true
			break
		}
	}

	return exists, index
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
	customerExists, _ := customer.findById()
	if customerExists {
		w.WriteHeader(http.StatusConflict)
		response, err := json.Marshal(struct {
			Status string `json:"status"`
		}{
			Status: "Customer with this identification already exists",
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
			break
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

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, CONTENT_VALUE)

	var customer Customer
	reqBody, _ := io.ReadAll(r.Body)

	json.Unmarshal(reqBody, &customer)

	customerExists, index := customer.findById()
	if customerExists {
		customers[index] = customer
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(customer)
	} else {
		w.WriteHeader(http.StatusNotFound)
		response, err := json.Marshal(struct {
			Status string `json:"status"`
		}{
			Status: fmt.Sprintf("Customer with id %d does not exists", customer.ID),
		})
		if err != nil {
			log.Fatal(err)
		}

		var message map[string]string

		json.Unmarshal(response, &message)
		json.NewEncoder(w).Encode(message)
	}
}
