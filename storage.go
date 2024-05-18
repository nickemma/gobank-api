package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Storage is an interface for storing and retrieving account data
type Storage interface {
	CreateAccount(account *Account) error
	DeleteAccount(id int) error
	UpdateAccount(account *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
}

// PostgresStore is a PostgreSQL implementation of the Storage interface
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgresStore instance
func NewPostgresStore() (*PostgresStore, error) {
	// Define the connection string for connecting to the Postgres database server
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	// Check if there was an error connecting to the database server
	if err != nil {
		return nil, err
	}

	// Check if the connection to the database server was successful
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return a new PostgresStore instance with the connected database
	return &PostgresStore{db: db}, nil
}

// Init initializes the PostgresStore by creating the account table
func (s *PostgresStore) Init() error {
	// Create the account table in the Postgres database
	return s.CreateAccountTable()
}

// CreateAccountTable creates the account table in the Postgres database
func (s *PostgresStore) CreateAccountTable() error {
	// Define the SQL query for creating the account table
	query := `CREATE TABLE IF NOT EXISTS account (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        account_number SERIAL,
        balance SERIAL,
        created_at TIMESTAMP
    );`

	// Execute the SQL query to create the account table
	_, err := s.db.Exec(query)
	return err
}

// CreateAccount inserts a new account into the account table
func (s *PostgresStore) CreateAccount(account *Account) error {
	// Define the SQL query for inserting a new account
	query := `INSERT INTO account (first_name, last_name, account_number, balance, created_at) 
               VALUES ($1, $2, $3, $4, $5)`

	// Execute the SQL query with the provided account data
	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.AccountNumber, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`DELETE FROM account WHERE id = $1`, id)
	return err
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id = $1`, id)

	// Check if there was an error executing the SQL query
	if err != nil {
		return nil, err
	}

	// Iterate over the rows returned by the SQL query
	for rows.Next() {
		return ScanInToAccount(rows)
	}
	// Return the account returned by the SQL query
	return nil, fmt.Errorf("account not found: %v", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account`)

	// Check if there was an error executing the SQL query
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the accounts
	var accounts []*Account

	// Iterate over the rows returned by the SQL query
	for rows.Next() {
		account, err := ScanInToAccount(rows)
		// Check if there was an error scanning the values
		if err != nil {
			return nil, err
		}

		// Append the account to the slice of accounts
		accounts = append(accounts, account)
	}

	// returned by the SQL query
	return accounts, nil
}

func ScanInToAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	// Scan the values from the current row into the account struct
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.AccountNumber,
		&account.Balance,
		&account.CreatedAt,
	)

	// Check if there was an error scanning the values
	if err != nil {
		return nil, err
	}

	return account, nil
}
