package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getStocks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	s := readLog()
	io.WriteString(w, s)
}

func readLog() string {
	data, err := ioutil.ReadFile("stocks.json")

	if err != nil {
		fmt.Println(err)
	}

	return string(data)
}