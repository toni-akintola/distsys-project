package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func createByteSlice(data any) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(p)

	if err != nil {
		fmt.Println("gob.Encode failed:", err)
	}

	return buf.Bytes()

}