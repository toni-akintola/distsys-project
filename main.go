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
	http.HandleFunc("/user/", accountHandler)
	marketServer := initializeMarket()
	fmt.Println(marketServer.data)
	// I'm working on the student machines and this is in the range of ports that work LOL
	err := http.ListenAndServe(":9444", nil)
	if err != nil {
		panic(err)
	}
		
}
