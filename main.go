package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set the JWT secret from the environment variable
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	store, err := newPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Drop existing tables
	if err := dropTables(store.db); err != nil {
		log.Fatal("Error dropping tables:", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Create a user for the Monopoly Bank
	monopolyUser := &User{
		FirstName: "Monopoly",
		LastName:  "Bank",
		Email:     "admin@gmail.com", // Set the email
		Password:  "megapassword",    // Set the password (ensure to hash it in production)
		CreatedAt: time.Now().UTC(),
	}

	// Create the user in the database
	if err := store.CreateUser(monopolyUser); err != nil {
		log.Fatal("Error creating Monopoly Bank user:", err)
	}

	// Create the Monopoly Bank account with the user ID
	if err := createMonopolyBankAccount(store, monopolyUser.ID); err != nil {
		log.Printf("Error creating Monopoly Bank account: %v", err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}

// Function to create the Monopoly Bank account
func createMonopolyBankAccount(store Storage, userID int) error {
	// Check if the Monopoly Bank account already exists
	existingAccount, err := store.GetAccountByID(1) // Change to the correct ID if needed
	if err != nil {
		return err
	}

	if existingAccount != nil {
		log.Println("Monopoly Bank account already exists.")
		return nil // Account already exists, no need to create
	}

	// Create a new Monopoly Bank account with the user ID
	monopolyAccount := &Account{
		FirstName: "Monopoly",
		LastName:  "Bank",
		Number:    999999,
		Balance:   999999999,
		CreatedAt: time.Now().UTC(),
		UserID:    &userID, // Assign the user ID here
	}

	return store.CreateAccount(monopolyAccount)
}

func dropTables(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS account CASCADE;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS users CASCADE;`)
	return err
}
