package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
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
	router.HandleFunc("/login", makeHTTPHandler(s.handleLogin)).Methods("POST")
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount)).Methods("GET")
	router.HandleFunc("/account/create", makeHTTPHandler(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/transfer", makeHTTPHandler(s.handleTransfer)).Methods("POST")
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandler(s.handleGetAccountByID), s.store))

	// Start the HTTP server
	log.Println("JSON API server is running on", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// handleLogin handles requests to the /login endpoint
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("unsupported method %s", r.Method)
	}
	// Create a new LoginRequest struct
	var loginReq LoginRequest

	// Decode the JSON request body into the LoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return err
	}

	// Close the request body
	defer r.Body.Close()

	// Get the account from the database
	acc, err := s.store.GetAccountByNumber(int(loginReq.Number))

	if err != nil {
		return err
	}

	fmt.Println(acc)

	return WriteJSON(w, http.StatusOK, loginReq)
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

	return fmt.Errorf("unsupported method %s", r.Method)
}

// handleGetAccount handles requests to the /accounts endpoint
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// handleGetAccount handles requests to the /account/{id} endpoint
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("unsupported method %s", r.Method)
}

// handleDeleteAccount handles requests to the /account/{id} endpoint
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"account deleted": id})
}

// handleCreateAccount handles requests to the /account endpoint
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// Create a new CreateAccountRequest struct
	createAccountReq := new(CreateAccountRequest)

	// Decode the JSON request body into the CreateAccountRequest struct
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	// Close the request body
	defer r.Body.Close()

	// Create a new account with the provided first and last name
	account, err := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)

	if err != nil {
		return err
	}

	// Store the account in the database
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	// Return the account as JSON
	return WriteJSON(w, http.StatusOK, account)
}

// handleTransfer handles requests to the /transfer endpoint
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// Create a new CreateAccountRequest struct
	transferReq := new(TransferAccountRequest)

	// Decode the JSON request body into the CreateAccountRequest struct
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	// Close the request body
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
	// Transfer the money between the accounts
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

// getID extracts the ID from the request URL
func getID(r *http.Request) (int, error) {
	idx := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idx)

	if err != nil {
		return id, fmt.Errorf("invalid account ID: %v", idx)
	}

	return id, nil
}

// withJWTAuth is a middleware that authenticates requests using JWT
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling the auth middleware")

		// Get the Authorization header
		tokenString := r.Header.Get("x-auth-token")

		// Validate the JWT token
		token, err := validateJWT(tokenString)

		// Check for errors
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Access denied"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			return
		}

		// Get the user ID from the token
		userID, err := getID(r)
		// Check for errors in getting the user ID
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
			return
		}
		// Get the account from the database
		account, err := s.GetAccountByID(userID)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
			return
		}

		// Get the claims from the token
		claims := token.Claims.(jwt.MapClaims)

		// Check if the account number in the token matches the account number in the database
		if account.AccountNumber != int64(claims["account_number"].(float64)) {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
			return
		}

		// Call the original handler function
		handlerFunc(w, r)
	}
}

// validateJWT validates a JWT token
func validateJWT(tokenString string) (*jwt.Token, error) {
	// Get the JWT secret from the environment
	secret := os.Getenv("JWT_SECRET")

	// Parse the token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(secret), nil
	})

}

// createJWT creates a new JWT token for the account
func createJWT(account *Account) (string, error) {
	// Get the JWT secret from the environment
	secret := os.Getenv("JWT_SECRET")

	// Set the claims
	claims := &jwt.MapClaims{
		"expiresAt":      time.Now().Add(time.Hour * 24).Unix(),
		"account_number": account.AccountNumber,
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the token string
	return token.SignedString([]byte(secret))
}
