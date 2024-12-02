package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type MarketServer struct {
	data map[string]*Stock
}

type Stock struct {
	Ticker string
	CompanyName string
	AssetType string
	CurrentPrice float64
	Currency string
	Volume float64
	LastUpdated string
	Volatility float64
}

func (s *MarketServer) readLog() {
	data, err := ioutil.ReadFile("stocks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var rawData map[string][]map[string]interface{}
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, rawStock := range rawData["stocks"] {

		// Extract and parse fields
		ticker, _ := rawStock["ticker"].(string)
		companyName, _ := rawStock["companyName"].(string)
		assetType, _ := rawStock["assetType"].(string)
		currency, _ := rawStock["currency"].(string)
		lastUpdated, _ := rawStock["lastUpdated"].(string)

		// Parse numeric fields
		currentPrice := 0.0
		volume := 0.0
		volatility := 0.0

		if cp, ok := rawStock["currentPrice"].(float64); ok {
			currentPrice = cp
		}

		if vol, ok := rawStock["volume"].(float64); ok {
			volume = vol
		}

		if vola, ok := rawStock["volatility"].(float64); ok {
			volatility = vola
		}


		// Create Stock Struct
		stock := &Stock{
			Ticker:       ticker,
			CompanyName:  companyName,
			AssetType:    assetType,
			CurrentPrice: currentPrice,
			Currency:     currency,
			LastUpdated:  lastUpdated,
			Volume:       volume,
			Volatility:   volatility,
		}

		s.data[ticker] = stock
		fmt.Println("Parsed stock:", stock)
	}
}

func (s *MarketServer) handleGetStock(w http.ResponseWriter, r *http.Request) {
	ticker := strings.TrimPrefix(r.URL.Path, "/single-stock/")
	if ticker == "" {
		http.Error(w, "please enter username", http.StatusBadRequest)
		return
	}
	fmt.Println(ticker)

	stock := s.getStock(ticker)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func (s *MarketServer) handleGetAllStocks(w http.ResponseWriter, r *http.Request) {
	ticker := strings.TrimPrefix(r.URL.Path, "/single-stock/")
	if ticker == "" {
		http.Error(w, "please enter username", http.StatusBadRequest)
		return
	}

	stock := s.getStock(ticker)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func (s MarketServer) getStock(ticker string) Stock {
	return *s.data[ticker]
}


func (s *MarketServer) updateStock(ticker string, volume float64, newPrice float64) {
	if stock, ok := s.data[ticker]; ok {
		stock.Volume = volume
		stock.CurrentPrice = newPrice
	}
}


func (s *MarketServer) writeLog() {
	// Open or create the log file
	file, err := os.OpenFile("market_log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	logger := log.New(file, "LOG: ", log.Ldate|log.Ltime)

	for ticker, stock := range s.data {
		logger.Printf("Ticker: %s, Current Price: %.2f\n", ticker, stock.CurrentPrice)
	}

	fmt.Println("Log written successfully.")
}

func (s *MarketServer) initializeMarket() {
	s.data = make(map[string]*Stock, 0);
	s.readLog();
}
 