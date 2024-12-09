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
	marketHost := "http://localhost:9444"
	executorServer, _ := newExecutorServer(marketHost)
	mux1, mux2 := http.NewServeMux(), http.NewServeMux()
	
	http.HandleFunc("/", getRoot)
	mux1.HandleFunc("/all-stocks/", marketServer.handleGetAllStocks)
	mux1.HandleFunc("/single-stock/", marketServer.handleGetStock)
	mux1.HandleFunc("/order/", marketServer.handleOrder)
	mux2.HandleFunc("/user/", executorServer.accountHandler)
	mux2.HandleFunc("/create-account/", executorServer.handleCreateAccount)
	mux2.HandleFunc("/single-stock/", executorServer.handleGetStock)
	mux2.HandleFunc("/all-stocks/", executorServer.handleGetAllStocks)
	mux2.HandleFunc("/order/", executorServer.handleOrder)
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

	time.Sleep(2 * time.Second)
	const numWorkers int = 5
	
	// Set up a timer so we can write to the log every 60 seconds
	logTicker := time.NewTicker(60 * time.Second)
	defer logTicker.Stop()

	updateTicker := time.NewTicker(2 * time.Second)
	defer updateTicker.Stop()

	go func (){
		for range updateTicker.C {
			for i := 0; i < numWorkers; i++ {
			fmt.Println("Random updating.")
				marketServer.randomUpdate()
				marketServer.writeLog()
				executorServer.saveAccounts()
			}
		}
			
	}()
	

	// Block the main goroutine
	select {}

		
}
