package main

import (
	"math/rand"
	"time"
)

// CreateAccountRequest is a request to create a new account
type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Account is a bank account model
type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountNumber int64     `json:"account_number"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

// NewAccount creates a new account with the given first and last name
func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: int64(rand.Intn(10000000)),
		CreatedAt:     time.Now().UTC(),
	}
}
