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
	server, _ := newAccountServer()
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/stocks/", marketServer.handleGetStock)
	http.HandleFunc("/user/", server.accountHandler)
	fmt.Println(marketServer.data)
	// Set up a ticker to run every 60 seconds
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Run the logging in a separate goroutine
	go func() {
		for range ticker.C {
			marketServer.writeLog()
		}
	}()
	// I'm working on the student machines and this is in the range of ports that work LOL
	marketErr := http.ListenAndServe(":9444", nil)
	executorErr := http.ListenAndServe(":9445", nil)
	if marketErr != nil {
		panic(marketErr)
	} else if executorErr != nil {
		panic(executorErr)
	}
	fmt.Println("Market server is running.")
	fmt.Println("Executor server is running.")
		
}
