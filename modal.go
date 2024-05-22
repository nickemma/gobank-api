package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// loginRequest is a request to login
type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

// Transfer is a request to transfer money between accounts
type TransferAccountRequest struct {
	ToAccountNumber int64 `json:"account_number"`
	Amount          int64 `json:"amount"`
}

// UnmarshalJSON unmarshals the JSON data into the TransferAccountRequest struct
func (s *TransferAccountRequest) UnmarshalJSON(data []byte) error {
	// Temporary struct for known fields
	type Alias TransferAccountRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Unmarshal into map to check for unknown fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Unmarshal into the known struct fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Check for unknown fields
	for key := range raw {
		switch key {
		case "account_number", "amount":
			// known fields
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}
	return nil
}

// CreateAccountRequest is a request to create a new account
type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// UnmarshalJSON unmarshals the JSON data into the CreateAccountRequest struct
func (s *CreateAccountRequest) UnmarshalJSON(data []byte) error {
	// Temporary struct for known fields
	type Alias CreateAccountRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Unmarshal into map to check for unknown fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Unmarshal into the known struct fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Check for unknown fields
	for key := range raw {
		switch key {
		case "first_name", "last_name", "password":
			// known fields
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}
	return nil
}

// Account is a bank account model
type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"_"`
	AccountNumber     int64     `json:"account_number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

// NewAccount creates a new account with the given first and last name
func NewAccount(firstName, lastName, password string) (*Account, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encrypt),
		AccountNumber:     int64(rand.Intn(10000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
