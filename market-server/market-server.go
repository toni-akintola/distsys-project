package marketserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MarketServer struct {
	data map[string]interface{}
}


func readLog() map[string]interface{} {
	data, err := ioutil.ReadFile("stocks.json")
	var result map[string]interface{}
	if err != nil {
		fmt.Println(err)
	}
	
	json.Unmarshal([]byte(data), &result)
	fmt.Println(result)
	return result
}

func initialize(s MarketServer) {
	s.data = readLog()
}