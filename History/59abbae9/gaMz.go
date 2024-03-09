package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	appName       = "5P50603710"
	appSource     = "10427"
	userID        = "uxuZEFys5nv"
	password      = "7elTHyW0EC3"
	userKey       = "sR12m8nkT8VEPXtfgLFlspj5BQlSqB51"
	encryptionKey = "jTS6yEtvhXThvDTYNHQNVXmklWFEaeQj"
)

func main() {

	// Access Token and Client Code
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6IjUwNjAzNzEwIiwicm9sZSI6InNSMTJtOG5rVDhWRVBYdGZnTEZsc3BqNUJRbFNxQjUxIiwiU3RhdGUiOiIiLCJuYmYiOjE2OTI5MzkyODYsImV4cCI6MTY5MjkzOTM0NiwiaWF0IjoxNjkyOTM5Mjg2fQ.AWX2Hm4oGRl-iyx7cYaWGd8D5kPUGEc9RhHFnLUtRUY"
	clientCode := "50603710"

	// WebSocket URL
	socketURL := fmt.Sprintf("wss://openfeed.5paisa.com/Feeds/api/chat?Value1=%s|%s", accessToken, clientCode)

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
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
