package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type MarketServer struct {
	data map[string]Stock
}

type Stock struct {
	ticker string
	companyName string
	assetType string
	currentPrice float64
	currency string
	volume float64
	lastUpdated string
	volatility float64
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
	// Create and initialize the map
	result := make(map[string]Stock) 

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
			volume = vola
		}


		// Create Stock Struct
		stock := Stock{
			ticker:       ticker,
			companyName:  companyName,
			assetType:    assetType,
			currentPrice: currentPrice,
			currency:     currency,
			lastUpdated:  lastUpdated,
			volume:       volume,
			volatility:   volatility,
		}

		result[ticker] = stock
		fmt.Println("Parsed stock:", stock)
	}
}

func (s *MarketServer) handleGetStock(w http.ResponseWriter, r *http.Request) {
	var result map[string]Stock
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
	}

	var tickerStock Stock = result["ticker"]

	fmt.Println(result, s.getStock(tickerStock.ticker))
}

func (s MarketServer) getStock(ticker string) Stock {
	return s.data[ticker]
}



func (s *MarketServer) initializeMarket() {
	s.readLog();

}
