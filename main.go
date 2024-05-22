package main

import (
	"flag"
	"fmt"
	"log"
)

// seedUsersAccount seeds the database with user accounts
func seedUsersAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)

	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created account: %d\n", acc.AccountNumber)
	return acc
}

func seedUsers(store Storage) {
	seedUsersAccount(store, "Techie", "Emma", "password22")
}

// main is the entry point for the application
func main() {
	// Define a flag for seeding the database with user accounts
	seed := flag.Bool("seed", false, "seed the database with user accounts")

	// Parse the command line flags
	flag.Parse()

	// Create a new PostgresStore instance
	store, err := NewPostgresStore()

	// Check if there was an error creating the PostgresStore instance
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the PostgresStore
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding the database with user accounts...")
		// Seed the database with user accounts
		seedUsers(store)
	}

	// Create a new APIServer instance
	server := NewAPIServer(":5000", store)

	// Start the APIServer
	server.Run()
}
