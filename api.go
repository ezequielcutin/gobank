package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
	var transferReq struct {
		FromID int64 `json:"fromId"`
		ToID   int64 `json:"toId"`
		Amount int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		return err
	}

	// Implement the transfer logic
	if err := s.store.TransferFunds(transferReq.FromID, transferReq.ToID, transferReq.Amount); err != nil {
		return err
	}

	// Create a response message
	responseMessage := fmt.Sprintf("Transaction of %d amount from ID %d to ID %d", transferReq.Amount, transferReq.FromID, transferReq.ToID)

	return WriteJSON(w, http.StatusOK, map[string]string{"message": responseMessage})
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
