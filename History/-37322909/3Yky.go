package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Replace with your Socket.IO server URL
	serverURL := "http://localhost:3000/socket.io/"

	// Establish a WebSocket connection to the Socket.IO server
	socket, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
	}
	defer socket.Close()

	// Handle incoming messages
	go func() {
		for {
			_, message, err := socket.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			// Parse the received message as per Socket.IO's protocol
			messageType, data, err := parseSocketIOMessage(message)
			if err != nil {
				log.Println("Error parsing Socket.IO message:", err)
				return
			}

			if messageType == "message" {
				fmt.Println("Received message:", string(data))
			} else if messageType == "event" {
				// Handle custom events here
				fmt.Println("Received custom event:", string(data))
			}
		}
	}()

	// Send a sample message to the server
	sampleMessage := `42["my-custom-event",{"message":"Hello, Socket.IO!"}]`
	err = socket.WriteMessage(websocket.TextMessage, []byte(sampleMessage))
	if err != nil {
		log.Println("Error sending message:", err)
	}

	// Handle interrupt signal to gracefully close the connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-interrupt:
			fmt.Println("Interrupt signal received. Closing WebSocket.")

			// Gracefully close the WebSocket connection
			socket.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))

			return
		}
	}
}

// Parse the incoming message based on Socket.IO's protocol
func parseSocketIOMessage(message []byte) (string, []byte, error) {
	// Implement the parsing logic here based on Socket.IO's protocol
	// Extract the message type and data

	// For example, assuming the message is in the format "42[<type>,<data>]"
	// Parse the message and return the message type and data

	// Replace this example logic with the actual Socket.IO parsing code
	messageType := "message"
	data := message

	return messageType, data, nil
}
