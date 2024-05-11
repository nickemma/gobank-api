package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// APIServer is an HTTP server that exposes a JSON API
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates a new APIServer
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the APIServer
func (s *APIServer) Run() {
	// Create a new router instance
	router := mux.NewRouter()

	// Define the routes for the API
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccountByID))

	// Start the HTTP server
	log.Println("JSON API server is running on", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// handleAccount handles requests to the /account endpoint
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Check the HTTP method and call the appropriate handler function
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	// Return an error for unsupported HTTP methods
	return fmt.Errorf("unsupported method %s", r.Method)
}

// handleGetAccount handles requests to the /accounts endpoint
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handleGetAccount handles requests to the /account/{id} endpoint
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println("ID:", id)

	return WriteJSON(w, http.StatusOK, &Account{})
}

// handleDeleteAccount handles requests to the /account/{id} endpoint
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handleCreateAccount handles requests to the /account endpoint
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// Create a new CreateAccountRequest struct
	createAccountReq := new(CreateAccountRequest)

	// Decode the JSON request body into the CreateAccountRequest struct
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	// Create a new account with the provided first and last name
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	// Store the account in the database
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	// Return the account as JSON
	return WriteJSON(w, http.StatusOK, account)
}

// handleTransfer handles requests to the /transfer endpoint
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ================= Helper Functions =================

// writeJSON writes the JSON representation of v to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode the value v as JSON and write it to the http.ResponseWriter
	return json.NewEncoder(w).Encode(v)
}

// apiFunc is a function that handles an API request
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// ApiError is an error response from the API
type ApiError struct {
	Error string `json:"error"`
}

// makeHTTPHandler creates an http.HandlerFunc from an apiFunc
func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	// Return an http.HandlerFunc that calls the apiFunc and handles any errors
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
