package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUsers() ([]*User, error)
	GetUserByID(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	TransferFunds(fromID int64, toID int64, amount int64) error
	GetBalance(id int) (int64, error)
	GetTransactions(id int) ([]Transaction, error)
}

type PostgresStore struct {
	db *sql.DB
}

func newPostGresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.createUsersTable(); err != nil {
		return err
	}
	if err := s.createTransactionsTable(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) createUsersTable() error {
	query := `create table if not exists users (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(100) unique not null,
		password varchar(100) not null,
		created_at timestamp,
		balance BIGINT DEFAULT 100
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createTransactionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        amount INTEGER NOT NULL,
        type VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
	log.Printf("Original (already hashed) password: %s", user.Password)

	log.Printf("Creating user with email: %s", user.Email)
	log.Printf("Hashed password to be stored: %s", user.Password)

	query := `INSERT INTO users (first_name, last_name, email, password, created_at, balance) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt, user.Balance)
	return err
}

func (s *PostgresStore) DeleteUser(id int) error {
	query := `delete from users where id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user with ID %d: %v", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user ID %d: %v", id, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No user found with ID %d", id)
		return fmt.Errorf("no user found with ID %d", id)
	}

	return nil
}

func (s *PostgresStore) UpdateUser(user *User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, balance = $5 WHERE id = $6`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Balance, user.ID)
	return err
}

func (s *PostgresStore) GetUserByID(id int) (*User, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at, balance FROM users WHERE id = $1`
	var user User
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Other error
	}
	return &user, nil
}

func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("SELECT * FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	users := []*User{}
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *PostgresStore) TransferFunds(fromID, toID int64, amount int64) error {
	log.Printf("Starting transfer: From ID %d to ID %d, Amount: %d", fromID, toID, amount)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	var fromBalance int64
	err = tx.QueryRow(`SELECT balance FROM users WHERE id = $1`, fromID).Scan(&fromBalance)
	if err != nil {
		log.Printf("Error fetching sender balance: %v", err)
		return err
	}
	log.Printf("Sender (ID: %d) balance: %d", fromID, fromBalance)

	if fromBalance < amount {
		log.Printf("Insufficient funds: Balance %d, Amount %d", fromBalance, amount)
		return fmt.Errorf("insufficient funds in account ID %d", fromID)
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		log.Printf("Error updating sender balance: %v", err)
		return err
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		log.Printf("Error updating recipient balance: %v", err)
		return err
	}

	//insert transaction records:
	_, err = tx.Exec(`INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, $3)`, fromID, -amount, "Sent")
	if err != nil {
		log.Printf("Error inserting sender transaction: %v", err)
		return err
	}

	_, err = tx.Exec(`INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, $3)`, toID, amount, "Received")
	if err != nil {
		log.Printf("Error inserting recipient transaction: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Println("Transfer completed successfully")
	return nil
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email = $1`

	var user User
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	log.Printf("Retrieved user from database - Email: %s, Hashed Password: %s", user.Email, user.Password)

	return &user, nil
}

func (s *PostgresStore) GetBalance(id int) (int64, error) {
	var balance int64
	err := s.db.QueryRow("SELECT balance FROM users WHERE id = $1", id).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *PostgresStore) GetTransactions(id int) ([]Transaction, error) {
	rows, err := s.db.Query("SELECT id, amount, type, created_at FROM transactions WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
