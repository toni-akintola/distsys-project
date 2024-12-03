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

// func runMarketServer() {

// }

// func runExecutorServer() {

// }

func main() {
	marketServer := &MarketServer{}
	marketServer.initializeMarket()
	marketHost := "http://localhost:9444"
	executorServer, _ := newExecutorServer(marketHost)
	mux1, mux2 := http.NewServeMux(), http.NewServeMux()
	
	http.HandleFunc("/", getRoot)
	mux1.HandleFunc("/all-stocks/", marketServer.handleGetAllStocks)
	mux1.HandleFunc("/single-stock/", marketServer.handleGetStock)
	mux2.HandleFunc("/user/", executorServer.accountHandler)
	mux2.HandleFunc("/single-stock/", executorServer.handleGetStock)
	mux2.HandleFunc("/all-stocks/", executorServer.handleGetAllStocks)
	// Create the first server
	server1 := &http.Server{
		Addr:    ":9444",
		Handler: mux1,
	}

	// Create the second server
	server2 := &http.Server{
		Addr:    ":9445",
		Handler: mux2,
	}

	// Start the servers in separate goroutines
	go func() {
		// I'm working on the student machines and this is in the range of ports that work LOL
		fmt.Println("Starting market server on port 9444.")
		if err := server1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Market server failed: %v\n", err)
		}
	}()

	go func() {
		fmt.Println("Starting executor server on port 9445.")
		if err := server2.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Executor server failed: %v\n", err)
		}
	}()

	// Set up a timer so we can write to the log every 60 seconds
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Run the logging in a separate goroutine
	go func() {
		for range ticker.C {
			marketServer.writeLog()
		}
	}()

	// Block the main goroutine
	select {}

	

		
}
