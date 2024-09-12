package main

import (
	"log"
	"time"
)

func main() {
	store, err := newPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Create the Monopoly Bank account
	if err := createMonopolyBankAccount(store); err != nil {
		log.Printf("Error creating Monopoly Bank account: %v", err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}

// Create the Monopoly Bank account, important for overdrafting
func createMonopolyBankAccount(store Storage) error {
	// Check if the Monopoly Bank account already exists
	existingAccount, err := store.GetAccountByID(21) // Change to whichever Monopoly Bank account ID you want to use
	if err != nil {
		return err
	}

	if existingAccount != nil {
		log.Println("Monopoly Bank account already exists.")
		return nil // Account already exists, no need to create
	}

	monopolyAccount := &Account{
		FirstName: "Monopoly",
		LastName:  "Bank",
		Number:    999999,
		Balance:   999999999,
		CreatedAt: time.Now().UTC(),
	}

	return store.CreateAccount(monopolyAccount)
}
