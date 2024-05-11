package main

import "database/sql"

type Storage interface {
	CreateAccount(account *Account) error
	DeleteAccount(id int) error
	UpdateAccount(account *Account) error
	DeleteAccountByID(id int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}
