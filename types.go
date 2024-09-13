package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// type Account struct {
// 	ID        int       `json:"id"`
// 	FirstName string    `json:"firstName"`
// 	LastName  string    `json:"lastName"`
// 	Number    int64     `json:"number"`
// 	Balance   int64     `json:"balance"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UserID    *int      `json:"userId"` // Added UserID to link to User
// }

// func NewAccount(firstName, LastName string) *Account {
// 	return &Account{
// 		FirstName: firstName,
// 		LastName:  LastName,
// 		Number:    int64(rand.Intn(1000000)),
// 		CreatedAt: time.Now().UTC(),
// 	}
// }

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // The "-" means this field won't be included in JSON output
	CreatedAt time.Time `json:"createdAt"`
	Balance   int64     `json:"balance"` // This can represent the account balance
	Number    int64     `json:"number"`  // This can represent the account number
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
	}, nil
}

type Transaction struct {
	ID        int       `json:"id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}
