package main

import (
	"fmt"
	"net/http"
)



func main() {
	fmt.Println("Server is running.")
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/stocks", getStocks)
	
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
		
}