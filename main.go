package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server is running.")
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/stocks", getStocks)
//	
	err := http.ListenAndServe(":9444", nil)
	if err != nil {
		panic(err)
	}
		
}
