package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)
type Account struct {
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
	positions []Position
}

type Order struct {
	Quantity float64 `json:"quantity"`
	Ticker string `json:"ticker"`
	Username string `json:"username"`
}

type Position struct {
	Order Order `json:"order"`
	Price float64 `json:"price"`
}


type ExecutorServer struct {
	data map[string]Account // Store accounts in a map for efficient lookups
	marketHost string
	httpClient http.Client
}

// NewExecutorServer initializes the ExecutorServer and loads accounts from a file
func newExecutorServer(marketHost string) (*ExecutorServer, error) {
	filename := "accounts.json"
	accounts, err := readAccounts(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %v", err)
	}

	data := make(map[string]Account)
	for _, account := range accounts {
		data[strings.ToLower(account.Username)] = account
	}
	executorClient := http.Client{
		Timeout: time.Second * 3, // Have our request timeout after 3 seconds if we don't receieve a response.
	}

	return &ExecutorServer{data: data, marketHost: marketHost, httpClient: executorClient}, nil
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

// saveAccounts writes account information to file after we update
func (s *ExecutorServer) saveAccounts () error {
	
	file, err := os.OpenFile("accounts.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open accounts.json for writing")
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s.data); err != nil {
		return fmt.Errorf("unable to encode accounts")
	}
	return nil
}

// handleCreateAccount creates a new account
func (s *ExecutorServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "please POST", http.StatusMethodNotAllowed)
		return
	}
	// parse JSON
	var newAccount Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newAccount); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// we should make sure that this username does not already exist
	if _, exists := s.data[strings.ToLower(newAccount.Username)]; exists {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}
	// give account $10k
	newAccount.Balance = 10000
	// store in-memory
	s.data[strings.ToLower(newAccount.Username)] = newAccount
	// store on-disk
	if err := s.saveAccounts(); err != nil {
		http.Error(w, fmt.Sprintf("unable to save account to disk"), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newAccount); err != nil {
		http.Error(w, "unable to encode JSON response", http.StatusInternalServerError)
	}
}

// getAccount retrieves an account by username
func (s *ExecutorServer) getAccount(username string) (*Account, bool) {
	account, exists := s.data[strings.ToLower(username)]
	if !exists {
		return nil, false
	}
	return &account, true
}

// updateAccount updates an account balance and position information for a given account.
// updateAccount updates an account's balance and position information.
func (s *ExecutorServer) updateAccount(username string, ticker string, quantity float64, price float64) bool {
	// Convert username to lowercase for consistent key lookup.
	accountKey := strings.ToLower(username)
	account, exists := s.data[accountKey]
	if !exists {
		return false // Account doesn't exist.
	}

	// Calculate the total cost of the transaction.
	totalCost := quantity * price
	if account.Balance < totalCost {
		return false // Insufficient funds.
	}

	// Deduct the cost from the account balance.
	account.Balance -= totalCost

	// Ensure positions slice is initialized if nil.
	if account.positions == nil {
		account.positions = make([]Position, 0)
	}

	// Create and append the new position.
	newPosition := Position{
		Order: Order{
			Ticker:   ticker,
			Quantity: quantity,
			Username: username,
		},
		Price: price,
	}
	account.positions = append(account.positions, newPosition)

	// Write the updated account back to the map.
	s.data[accountKey] = account

	return true
}

// AccountHandler handles HTTP requests for account information
func (s *ExecutorServer) accountHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *ExecutorServer) handleGetStock(w http.ResponseWriter, r *http.Request) {

	ticker := strings.TrimPrefix(r.URL.Path, "/single-stock/")
	if ticker == "" {
		http.Error(w, "please enter username", http.StatusBadRequest)
		return
	}
	
	response, responseError := http.Get(s.marketHost + "/single-stock/" + ticker)

	var result map[string]interface{}
	if responseError != nil {
		fmt.Println("Error sending request:", responseError)
		return
	}
	defer response.Body.Close()
	fmt.Println("Response status:", response.Status)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body.")
	}
	json.Unmarshal([]byte(body), &result)
	fmt.Println(result)

}

func (s *ExecutorServer) handleGetAllStocks(w http.ResponseWriter, r *http.Request) {
	response, responseError := http.Get(s.marketHost + "/all-stocks/")

	var result map[string]interface{}
	if responseError != nil {
		fmt.Println("Error sending request:", responseError)
		return
	}

	defer response.Body.Close()
	fmt.Println("Response status:", response.Status)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body.")
	}
	json.Unmarshal([]byte(body), &result)
	fmt.Println(result)

}

func (s *ExecutorServer) handleOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	err := json.NewDecoder(r.Body).Decode(&o)
	fmt.Println(o)
	// requestBody, err := unmarshalJSONBody[Order](r)
	

	if err != nil {
		fmt.Println("Error unmarshalling JSON body:", err)
		return
	}
	jsonBody := createRequestBody(o)
	response, responseError := http.Post(s.marketHost + "/order/", "application/json", jsonBody)

	
	if responseError != nil {
		fmt.Println("Error sending request:", responseError)
		return
	}

	defer response.Body.Close()
	fmt.Println("Response status:", response.Status)

	var p Position
	json.NewDecoder(response.Body).Decode(&p)
	s.updateAccount(p.Order.Username, p.Order.Ticker, p.Order.Quantity, p.Price)
	fmt.Println(s.data[p.Order.Username])

}
