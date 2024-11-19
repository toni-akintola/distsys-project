package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type MarketServer struct {
	data map[string]interface{}
}


func readLog() map[string]interface{} {
	data, err := ioutil.ReadFile("stocks.json")
	var result map[string]interface{}
	if err != nil {
		fmt.Println(err)
	}
	
	json.Unmarshal([]byte(data), &result)
	return result
}

func handleGetStock(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{}
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
func getStock(s MarketServer, ticker string) interface{} {
	return s.data[ticker]
}




func initializeMarket() *MarketServer {
	s := MarketServer {}
	s.data = readLog()

	return &s
}
