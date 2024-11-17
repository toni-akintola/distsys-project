package main

import (
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}


func main() {
	fmt.Println("Server is running.")
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/stocks", handleGetStock)
	marketServer := initializeMarket()
	fmt.Println(marketServer.data)
//	
	err := http.ListenAndServe(":9444", nil)
	if err != nil {
		panic(err)
	}
		
}
