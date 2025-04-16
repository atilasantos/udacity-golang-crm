package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	ContentType  = "Content-Type"
	ContentValue = "application/json"
)

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
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
		log.Fatal(fmt.Printf("Failed to load customers file: %s", err))
	}

	defer fileContent.Close()
	byteResult, _ := io.ReadAll(fileContent)

	var customers []Customer
	json.Unmarshal([]byte(byteResult), &customers)

	return customers

}

var customers []Customer = loadCustomers(os.Getenv("CUSTOMER_FILE_PATH"))

func AddCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentValue)

	var customer Customer
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
	}

	json.Unmarshal(reqBody, &customer)
	customerExists, _ := customer.findById()
	if customerExists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Customer with this identification already exists",
		})
	} else {
		customers = append(customers, customer)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customer)
	}
}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentValue)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentValue)

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
	w.Header().Set(ContentType, ContentValue)
	var found bool = false

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for index, customer := range customers {
		if customer.ID == id {
			found = true
			customers = slices.Delete(customers, index, index+1)
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(customer)
			return
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": fmt.Sprintf("Customer with id %d does not exists", id),
		})
	}
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentValue)

	vars := mux.Vars(r)
	pathIDStr := vars["id"]
	pathID, err := strconv.Atoi(pathIDStr)
	if err != nil {
		http.Error(w, "Invalid ID in path", http.StatusBadRequest)
		return
	}

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if customer.ID != pathID {
		http.Error(w, "ID in path and body do not match", http.StatusBadRequest)
		return
	}

	customerExists, index := customer.findById()
	if customerExists {
		customers[index] = customer
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(customer)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": fmt.Sprintf("Customer with id %d does not exist", customer.ID),
		})
	}
}
