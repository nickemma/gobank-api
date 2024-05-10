package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// writeJSON writes the JSON representation of v to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// APIServer is an HTTP server that exposes a JSON API
type APIServer struct {
	listenAddr string
}

// NewAPIServer creates a new APIServer
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// Run starts the APIServer
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccount))

	log.Println("JSON API server is running on", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// handleAccount handles requests to the /account endpoint
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("unsupported method %s", r.Method)
}

// handleGetAccount handles requests to the /account/{id} endpoint
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
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
	return nil
}

// handleTransfer handles requests to the /transfer endpoint
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
