package main

import (
	"net/http"
)



func main() {
	getAccountInformation()
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
		
}