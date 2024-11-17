package dataserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type DataServer struct {
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

func initialize(s DataServer) {
	s.data = readLog()
}