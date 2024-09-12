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

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer)).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", makeHTTPHandleFunc(s.handleRegister)).Methods("POST")
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed! %s", r.Method)
}

// GET /account
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", idStr)
	}

	account, err := s.store.GetAccountByID(id)
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

	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", idStr)
	}

	// Attempt to delete the account
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	// Return 204 No Content on successful deletion
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("missing Authorization header")
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := validateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// Extract user ID from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user ID in token")
	}

	// Use the user's ID as the fromID
	if err := s.store.TransferFunds(int64(userID), transferReq.ToID, transferReq.Amount); err != nil {
		return err
	}

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

	return WriteJSON(w, http.StatusCreated, user)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")                           // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS") // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")               // Allowed headers

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginReq := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(loginReq.Email)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return fmt.Errorf("incorrect password")
	}

	// Generate JWT token
	token, err := createJWT(user)
	if err != nil {
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
