package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Welcome to South Bend Bets!\n")
}

func updateNameServer(port int, name string) {
	// Define the message as a map
	message := map[string]interface{}{
		"type":    "south-bend-bets",
		"owner":   "oakintol",
		"port":    port,
		"project": name,
	}

	// Marshal the message into JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	// Resolve the address of the catalog server
	target, err := net.ResolveUDPAddr("udp", "catalog.cse.nd.edu:9097")
	if err != nil {
		panic(err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send the message
	_, err = conn.Write(messageBytes)
	if err != nil {
		panic(err)
	}
}

func main() {
		// First server
	listener1, err := net.Listen("tcp", ":0") // OS assigns a random port
	if err != nil {
		panic(err)
	}
	defer listener1.Close()

	marketPort := listener1.Addr().(*net.TCPAddr).Port
	fmt.Printf("Market server is listening on port %d\n", marketPort)

	// Second server
	listener2, err := net.Listen("tcp", ":0") // OS assigns a random port
	if err != nil {
		panic(err)
	}
	defer listener2.Close()
	executorPort := listener2.Addr().(*net.TCPAddr).Port

	marketServer := &MarketServer{}
	marketServer.initializeMarket()
	hostName, err := os.Hostname()
	if err != nil {
		return
	}
	marketHost := "http://" + hostName + ":" + strconv.Itoa(marketPort)
	executorServer, _ := newExecutorServer(marketHost)
	updateNameServer(executorPort, hostName)

	mux1, mux2 := http.NewServeMux(), http.NewServeMux()
	http.HandleFunc("/", getRoot)
	mux1.HandleFunc("/all-stocks", marketServer.handleGetAllStocks)
	mux1.HandleFunc("/single-stock", marketServer.handleGetStock)
	mux1.HandleFunc("/order", marketServer.handleOrder)
	mux2.HandleFunc("/account", executorServer.accountHandler)
	mux2.HandleFunc("/create-account", executorServer.handleCreateAccount)
	mux2.HandleFunc("/single-stock", executorServer.handleGetStock)
	mux2.HandleFunc("/all-stocks", executorServer.handleGetAllStocks)
	mux2.HandleFunc("/order", executorServer.handleOrder)



	server1 := &http.Server{
		Handler: mux1, // Replace with your mux1
	}
	go func() {
		if err := server1.Serve(listener1); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()


	fmt.Printf("Executor server is listening on port %d\n", executorPort)

	server2 := &http.Server{
		Handler: mux2, // Replace with your mux2
	}
	go func() {
		if err := server2.Serve(listener2); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()


	
	// Start the servers in separate goroutines
	go func() {
		// I'm working on the student machines and this is in the range of ports that work LOL
		if err := server1.Serve(listener1); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Market server failed: %v\n", err)
		}
	}()

	go func() {
		if err := server2.Serve(listener2); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Executor server failed: %v\n", err)
		}
	}()

	time.Sleep(2 * time.Second)
	const numWorkers int = 5
	
	// Set up a timer so we can update the name server every 60 seconds
	nameServerTicker := time.NewTicker(60 * time.Second)
	defer nameServerTicker.Stop()

	updateTicker := time.NewTicker(2 * time.Second)
	defer updateTicker.Stop()

	go func (){
		for range updateTicker.C {
			for i := 0; i < numWorkers; i++ {
				marketServer.randomUpdate()
				marketServer.writeLog()
				executorServer.saveAccounts()
			}
		}
			
	}()

	go func (){
		for range nameServerTicker.C {
			updateNameServer(executorPort, hostName)
		}
			
	}()
	

	// Block the main goroutine and select whichever goroutine is ready.
	select {}

		
}
