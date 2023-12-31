package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Type        string  `json:"type"`
	Remarks     string  `json:"remarks"`
	TotalAmount float64 `json:"totalamount"`
	Household   string  `json:"household"`
}

func main() {
	jsonStruct := Payload{
		Type:        "conservancy",
		Remarks:     "Test from go",
		TotalAmount: 123.45,
		Household:   "678",
	}

	payloadBytes, err := json.Marshal(jsonStruct)
	if err != nil {
		// handle err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:8080/payments/add", body)
	if err != nil {
		fmt.Println("new request err")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

}
