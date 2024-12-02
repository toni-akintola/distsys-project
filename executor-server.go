package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Position struct {
	quantity float64
	price float64
	ticker string
}
type Account struct {
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
	Positions []Position
}

type AccountServer struct {
	data map[string]Account // Store accounts in a map for efficient lookups
}

// NewAccountServer initializes the AccountServer and loads accounts from a file
func newAccountServer() (*AccountServer, error) {
	filename := "accounts.json"
	accounts, err := readAccounts(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %v", err)
	}

	data := make(map[string]Account)
	for _, account := range accounts {
		data[strings.ToLower(account.Username)] = account
	}

	return &AccountServer{data: data}, nil
}

// readAccounts reads accounts from a JSON file
func readAccounts(filename string) ([]Account, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	var accounts []Account
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&accounts); err != nil {
		return nil, fmt.Errorf("unable to decode JSON: %v", err)
	}

	return accounts, nil
}

// getAccount retrieves an account by username
func (s *AccountServer) getAccount(username string) (*Account, bool) {
	account, exists := s.data[strings.ToLower(username)]
	if !exists {
		return nil, false
	}
	return &account, true
}

// updateAccount updates an account balance and position information for a given account
func (s *AccountServer) updateAccount(username string, ticker string, quantity float64, price float64) bool {
	account, exists := s.data[strings.ToLower(username)]
	if !exists {
		return false
	}
	account.Balance -= price;
	position := Position{price: price, ticker: ticker, quantity: quantity}
	// If the positions list doesn't exist for an account, create it.
	if account.Positions == nil {
		account.Positions = make([]Position, 0)
	}
	account.Positions = append(account.Positions, position)
	return true;
}

// AccountHandler handles HTTP requests for account information
func (s *AccountServer) accountHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/user/")
	if username == "" {
		http.Error(w, "please enter username", http.StatusBadRequest)
		return
	}

	fmt.Printf("received request for username %s\n", username)

	account, found := s.getAccount(username)
	if !found {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		return
	}
}
