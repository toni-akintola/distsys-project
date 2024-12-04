package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func createByteSlice(data any) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)

	if err != nil {
		fmt.Println("gob.Encode failed:", err)
	}

	return buf.Bytes()

}

// Generic function to unmarshal JSON body into a given type
func unmarshalJSONBody[T any](r *http.Request) (T, error) {
	var result T
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return result, fmt.Errorf("failed to read body: %w", err)
	}
	

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func createRequestBody(data any) *bytes.Buffer {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Println("failed to read body: %w", err))
	}

	return bytes.NewBuffer(jsonBody)
}