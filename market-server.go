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
	volume int64
	lastUpdated string
	volatility float64
}

func readLog() map[string]Stock {
	data, err := ioutil.ReadFile("stocks.json")
	var rawData map[string][]map[string]interface{};
	var result map[string]Stock
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal([]byte(data), &rawData);

	for stock := range rawData["stocks"] {
		fmt.Println(rawData["stocks"][stock]);
	}

	return result
}

func (s *MarketServer) handleGetStock(w http.ResponseWriter, r *http.Request) {
	var result map[string]Stock
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
	}

	var tickerStock Stock = result["ticker"];

	fmt.Println(result, s.getStock(tickerStock.ticker));
}

func (s MarketServer) getStock(ticker string) Stock {
	return s.data[ticker];
}




func initializeMarket() *MarketServer {
	s := MarketServer {}
	s.data = readLog()

	return &s
}
