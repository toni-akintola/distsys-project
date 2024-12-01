package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}


func main() {
	marketServer := &MarketServer{}
	marketServer.initializeMarket()
	fmt.Println("Server is running.")
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/stocks", marketServer.handleGetStock)
	http.HandleFunc("/user/", accountHandler)
	fmt.Println(marketServer.data)
	// Set up a ticker to run every 60 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Run the logging in a separate goroutine
	go func() {
		for range ticker.C {
			marketServer.writeLog()
		}
	}()
	// I'm working on the student machines and this is in the range of ports that work LOL
	err := http.ListenAndServe(":9444", nil)
	if err != nil {
		panic(err)
	}
		
}
