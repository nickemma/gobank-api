package main

import "log"

func main() {
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

	// Create a new APIServer instance
	server := NewAPIServer(":5000", store)

	// Start the APIServer
	server.Run()
}
