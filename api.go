package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func createJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// fmt.Printf(token.SignedString(jwtSecret)) //token printing (debug)
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Add this struct at the top of the file with other request structs
type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCors(w, r)
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleDeleteUser)).Methods("DELETE")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetUserById))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer)).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", makeHTTPHandleFunc(s.handleRegister)).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST", "OPTIONS")
	router.HandleFunc("/balance/{id}", makeHTTPHandleFunc(s.handleGetBalance)).Methods("GET")
	router.HandleFunc("/transactions/{id}", makeHTTPHandleFunc(s.handleGetTransactions)).Methods("GET")
	router.HandleFunc("/user-by-email/{email}", makeHTTPHandleFunc(s.handleGetUserByEmail)).Methods("GET")
	router.HandleFunc("/user-details/{email}", makeHTTPHandleFunc(s.handleGetUserDetails)).Methods("GET")

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUser(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}

	return fmt.Errorf("method not allowed! %s", r.Method)
}

// GET /account
func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", idStr)
	}

	account, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("account not found with ID: %d", id)
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	//debugging purposes
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Println("Request Body:", string(body))

	//needed to include this to navigate EOF error
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	account, err := NewUser(createUserReq.FirstName, createUserReq.LastName, createUserReq.Email, createUserReq.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateUser(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %s", idStr)
	}

	// Attempt to delete the account
	if err := s.store.DeleteUser(id); err != nil {
		return err
	}

	// Return 204 No Content on successful deletion
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	log.Println("Starting transfer process")

	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		log.Printf("Error decoding transfer request: %v", err)
		return err
	}
	log.Printf("Transfer request decoded: %+v", transferReq)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Missing Authorization header")
		return fmt.Errorf("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := validateJWT(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		return fmt.Errorf("invalid token")
	}
	log.Println("JWT token validated successfully")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims")
		return fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Println("Invalid user ID in token")
		return fmt.Errorf("invalid user ID in token")
	}
	log.Printf("User ID extracted from token: %v", userID)

	err = s.store.TransferFunds(int64(userID), transferReq.ToID, transferReq.Amount)
	if err != nil {
		log.Printf("Error during transfer: %v", err)
		return err
	}

	log.Println("Transfer completed successfully")
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "Transfer successful"})
}

type TransferRequest struct {
	ToID   int64 `json:"toId"`
	Amount int64 `json:"amount"`
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	user, err := NewUser(createUserReq.FirstName, createUserReq.LastName, createUserReq.Email, createUserReq.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	token, err := createJWT(user)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token: token,
		User:  user,
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// CORS middleware
func enableCors(w http.ResponseWriter, r *http.Request) {
	log.Println("Received CORS request:", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")                            // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")  // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allowed headers

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginReq := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		log.Printf("Error decoding login request: %v", err)
		return err
	}

	user, err := s.store.GetUserByEmail(loginReq.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return fmt.Errorf("invalid credentials")
	}

	if user == nil {
		log.Printf("User not found for email: %s", loginReq.Email)
		return fmt.Errorf("invalid credentials")
	}

	log.Printf("Retrieved user from database - Email: %s, Hashed Password: %s", user.Email, user.Password)

	log.Printf("Comparing passwords for user: %s", loginReq.Email)
	log.Printf("Stored hashed password: %s", user.Password)
	log.Printf("Provided password: %s", loginReq.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return fmt.Errorf("incorrect password")
	}

	// Generate JWT token
	token, err := createJWT(user)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		return err
	}

	resp := LoginResponse{
		Token: token,
		User:  user,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func (s *APIServer) handleGetBalance(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %s", idStr)
	}

	balance, err := s.store.GetBalance(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int64{"balance": balance})
}

func (s *APIServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %s", idStr)
	}

	transactions, err := s.store.GetTransactions(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, transactions)
}

// func (s *APIServer) handleGetBalance(w http.ResponseWriter, r *http.Request) error {
//     id := mux.Vars(r)["id"]
//     // Convert id to int and fetch balance from database
//     // For now, let's return a dummy balance
//     return WriteJSON(w, http.StatusOK, map[string]int{"balance": 1000})
// }

// func (s *APIServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) error {
//     id := mux.Vars(r)["id"]
//     // Fetch transactions from database
//     // For now, let's return dummy transactions
//     transactions := []map[string]interface{}{
//         {"id": 1, "amount": 100, "type": "deposit"},
//         {"id": 2, "amount": 50, "type": "withdrawal"},
//     }
//     return WriteJSON(w, http.StatusOK, transactions)
// }

func (s *APIServer) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) error {
	email := mux.Vars(r)["email"]
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found with email: %s", email)
	}
	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleGetUserDetails(w http.ResponseWriter, r *http.Request) error {
	email := mux.Vars(r)["email"]
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return err
	}

	userDetails := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return WriteJSON(w, http.StatusOK, userDetails)
}
