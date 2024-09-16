package main

import (
	"database/sql"
	"log"

	// "os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
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

	monopolyPassword := "megapassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(monopolyPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing Monopoly Bank password:", err)
	}
	// Create a user for the Monopoly Bank
	monopolyUser := &User{
		FirstName: "Monopoly",
		LastName:  "Bank",
		Email:     "admin@gmail.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		Balance:   999999999, // Set initial balance
		Number:    999999,    // Set account number
	}

	// Create the user in the database
	if err := store.CreateUser(monopolyUser); err != nil {
		log.Fatal("Error creating Monopoly Bank user:", err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}

// Function to create the Monopoly Bank account
func createMonopolyBankAccount(store Storage, userID int) error {
	// Check if the Monopoly Bank account already exists
	existingAccount, err := store.GetUserByID(1) // Change to the correct ID if needed
	if err != nil {
		return err
	}

	if existingAccount != nil {
		log.Println("Monopoly Bank account already exists.")
		return nil // Account already exists, no need to create
	}

	// Create a new Monopoly Bank account with the user ID
	monopolyAccount := &User{
		FirstName: "Monopoly",
		LastName:  "Bank",
		Email:     "monopolybank@example.com", // Assuming a default email for the Monopoly Bank
		Password:  "monopolybankpassword",     // Assuming a default password for the Monopoly Bank
		CreatedAt: time.Now().UTC(),
		Balance:   999999999,
		Number:    999999,
	}

	return store.CreateUser(monopolyAccount)
}

func dropTables(db *sql.DB) error {
	// _, err := db.Exec(`DROP TABLE IF EXISTS account CASCADE;`)
	// if err != nil {
	// 	return err
	// }
	_, err := db.Exec(`DROP TABLE IF EXISTS users CASCADE;`)
	_, err = db.Exec(`DROP TABLE IF EXISTS transactions CASCADE;`)
	return err
}
