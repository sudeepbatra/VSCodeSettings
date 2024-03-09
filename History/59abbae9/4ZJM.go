package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// ... (constants and main function remain the same)

func main() {
	// ... (constants, access token, and client code remain the same)

	// Configure HTTP client
	httpClient := &http.Client{}

	// Create a request with the necessary headers
	request, err := http.NewRequest("GET", socketURL, nil)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}
	request.Header.Add("Sec-WebSocket-Protocol", "json")

	// Perform the request to initiate the WebSocket handshake
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatal("WebSocket handshake request error:", err)
	}
	defer resp.Body.Close()

	// Upgrade the response to a WebSocket connection
	conn, err := websocket.Upgrade(resp.Body, resp.Request, nil, 1024, 1024)
	if err != nil {
		log.Fatal("WebSocket upgrade error:", err)
	}
	defer conn.Close()

	// Start a goroutine to read messages from WebSocket
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				return
			}
			// Handle the received message, such as parsing the JSON and processing the data
			fmt.Println("Received:", string(message))
		}
	}()

	// Keep the main goroutine alive
	for {
		time.Sleep(time.Second) // You can perform other tasks here
	}
}
