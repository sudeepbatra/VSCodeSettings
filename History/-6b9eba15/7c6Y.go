package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	url := "https://api.icicidirect.com/breezeapi/api/v1/funds"
	payload := strings.NewReader("{}")

	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	correctCasingHeader := make(http.Header)
	correctCasingHeader["X-AppKey"] := ["s76162#+U35414Y*S413=099_FA6P567"]

	// Set the headers with the desired casing
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AppKey", "s76162#+U35414Y*S413=099_FA6P567")
	req.Header.Set("X-Checksum", "token 4ee90e3b5be8e7a3f4dbac60d806c73a01df1dea8f986ed4815c5f31f086c3e7")
	req.Header.Set("X-SessionToken", "U1VEREVQQkE6ODkyMTI0ODc=")
	req.Header.Set("X-Timestamp", "2023-07-29T03:30:58.000Z")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Process the response here...
}
