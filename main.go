package main

import (
	"fmt"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)



func main() {
	// Get our account information.

	client := alpaca.NewClient(alpaca.ClientOpts{
		// Alternatively you can set your key and secret using the
		// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
		APIKey:    os.Getenv("ALPACA_API_KEY"),
		APISecret: os.Getenv("ALPACA_SECRET"),
		BaseURL:   "https://paper-api.alpaca.markets",
	})
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)

	// Check if our account is restricted from trading.
	if acct.TradingBlocked {
		fmt.Println("Account is currently restricted from trading.")
	}

	// Check how much money we can use to open new positions.
	fmt.Printf("%v is available as buying power.\n", acct.BuyingPower)
}