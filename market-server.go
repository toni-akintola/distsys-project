package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
)
var TICKERS = []string{"AAPL", "TSLA", "AMZN", "JNG", "GOOGL"}
type MarketServer struct {
	data map[string]*Stock
}

type Stock struct {
	Ticker string `json:"ticker"`
	CompanyName string `json:"companyName"`
	AssetType string `json:"assetType"`
	CurrentPrice float64 `json:"currentPrice"`
	Currency string `json:"currency"`
	Volume float64 `json:"volume"`
	LastUpdated string `json:"lastUpdated"`
	Volatility float64 `json:"volatility"`
	SignTendency float64 `json:"signTendency"`
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
		signTendency := 0.0

		if cp, ok := rawStock["currentPrice"].(float64); ok {
			currentPrice = cp
		}

		if vol, ok := rawStock["volume"].(float64); ok {
			volume = vol
		}

		if vola, ok := rawStock["volatility"].(float64); ok {
			volatility = vola
		}

		if st, ok := rawStock["signTendency"].(float64); ok {
			signTendency = st
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
			SignTendency: signTendency,
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

	stock, err := s.getStock(ticker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func (s *MarketServer) handleGetAllStocks(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.data); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func (s MarketServer) getStock(ticker string) (*Stock, error) {
	stock, exists := s.data[ticker]
	// so that we do not dereference nonexistent pointer
	if !exists {
		return nil, fmt.Errorf("ticker not found")
	}
	return stock, nil
}

func (s *MarketServer) handleOrder(w http.ResponseWriter, r *http.Request) {

	// Unmarshal the JSON body into a Go struct
	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	
	var order Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	stock, _ := s.getStock(order.Ticker)
	position := Position{Order: order, Price: stock.CurrentPrice}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(position); err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
	
}

func (s *MarketServer) randomUpdate() {
	x := rand.IntN(len(TICKERS))
	ticker := TICKERS[x]
	stock, _ := s.getStock(ticker)
	signIndicator := rand.Float64()
	volatility := 0.0
	if signIndicator <= stock.SignTendency {
		volatility = stock.Volatility
	} else {
		volatility = -1 * stock.Volatility
	}

	dPrice := stock.CurrentPrice * volatility 
	fmt.Println("Ticker: ", ticker, "Old Price: ", stock.CurrentPrice, "New Price: ", stock.CurrentPrice + dPrice)
	// s.updateStock(ticker, stock.CurrentPrice + dPrice)
	
}


func (s *MarketServer) updateStock(ticker string, newPrice float64) {
	if stock, ok := s.data[ticker]; ok {
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
 
