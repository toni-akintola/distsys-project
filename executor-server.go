package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// was thinking that we store our accounts in a JSON file ?
// we can maybe get access to db8 which is a machine here with MySQL
// but this might be too much , it is tricky to interface website and database
type Account struct {
	Username string `json:"username"`
	Balance float64 `json:"balance"`
}



func readAccounts(filename string) ([]Account, error) {
	file, err := os.Open(filename)
	
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	// defer ensures close of file no matter result of function
	defer file.Close()

	var accounts []Account
	decoder := json.NewDecoder(file)
	// JSON --> Account slice
	err = decoder.Decode(&accounts)
	
	if err != nil {
		return nil, fmt.Errorf("unable to decode JSON: %v", err)
	}

	return accounts, nil
}

func getAccount(accounts []Account, username string) (*Account, bool) {
	for _, account := range accounts {
		// case insensitive
		if strings.EqualFold(account.Username, username) {
			return &account, true
		}
	}

	return nil, false
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	// get username from URL ?
	// like : /user/dthain
	username := strings.TrimPrefix(r.URL.Path, "/user/")
	if username == "" {
		// throw 400
		http.Error(w, "please enter username", http.StatusBadRequest)
		return
	}
	// check username
	fmt.Printf("received request for username %s\n", username)
	// pull all account data ... we should require authentication for this
	accounts, err := readAccounts ("accounts.json")

	if err != nil {
		// throw 500
		http.Error(w, "error reading account data", http.StatusInternalServerError)
		return
	}
	info, found := getAccount (accounts, username)
	if !found {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}
	// return account information
	// browser needs to know to expect JSON response
	w.Header().Set("Content-Type", "application/json")
	// info to do something with
	err = json.NewEncoder(w).Encode(found)
	// did we find our account ?
	if err != nil {
		// throw 500
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(info)
	// what does our account look like ?
	if err != nil {
		// throw 500
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}
