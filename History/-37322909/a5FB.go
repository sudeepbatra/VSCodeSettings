package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SocketIOMessage struct {
	Type int      `json:"type"`
	Data []string `json:"data"`
}

func main() {
	// Replace with your Socket.IO server URL
	serverURL := "http://localhost:3000/socket.io/"

	// Configure WebSocket options
	websocketDialer := websocket.DefaultDialer
	websocketDialer.ReadBufferSize = 1024
	websocketDialer.WriteBufferSize = 1024

	var mu sync.Mutex // Mutex to protect the connection state
	connected := false

	// Initialize reconnection parameters
	reconnectAttempts := 0
	maxReconnectAttempts := 10
	initialReconnectDelay := 5 * time.Second

	// Main loop for reconnection attempts with exponential backoff
	for {
		// Attempt to establish a WebSocket connection to the Socket.IO server
		socket, _, err := websocketDialer.Dial(serverURL, nil)
		if err != nil {
			log.Printf("WebSocket connection error: %v", err)

			if reconnectAttempts >= maxReconnectAttempts {
				log.Fatal("Exceeded maximum reconnection attempts.")
				return
			}

			// Calculate exponential backoff delay
			reconnectDelay := calculateReconnectDelay(initialReconnectDelay, reconnectAttempts)

			log.Printf("Attempting to reconnect in %v...", reconnectDelay)
			time.Sleep(reconnectDelay)
			reconnectAttempts++
			continue
		}
		log.Println("WebSocket connection established.")

		defer socket.Close()

		// Reset reconnection attempts upon successful connection
		reconnectAttempts = 0

		// Handle incoming messages
		go func() {
			for {
				_, message, err := socket.ReadMessage()
				if err != nil {
					log.Println("Error reading message:", err)
					return
				}

				messageType, data, err := parseSocketIOMessage(message)
				if err != nil {
					log.Println("Error parsing Socket.IO message:", err)
					return
				}

				switch messageType {
				case "connect":
					mu.Lock()
					connected = true
					mu.Unlock()
					fmt.Println("Connected to Socket.IO server")

				case "message":
					fmt.Println("Received message:", data)

				case "event":
					handleCustomEvent(data)

				case "disconnect":
					mu.Lock()
					connected = false
					mu.Unlock()
					fmt.Println("Disconnected from Socket.IO server")
					return // Exit the goroutine on disconnection
				}
			}
		}()

		// Send a sample message to the server
		sendSocketIOMessage(socket, "my-custom-event", map[string]string{"message": "Hello, Socket.IO!"})

		// Handle interrupt signal to gracefully close the connection
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		for {
			select {
			case <-interrupt:
				fmt.Println("Interrupt signal received. Closing WebSocket.")

				mu.Lock()
				if connected {
					socket.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
				}
				mu.Unlock()

				return
			}
		}
	}
}

func calculateReconnectDelay(initialDelay time.Duration, attempts int) time.Duration {
	// Calculate exponential backoff delay with a maximum limit
	maxDelay := 10 * time.Minute
	delay := initialDelay * time.Duration(math.Pow(2, float64(attempts)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

func parseSocketIOMessage(message []byte) (string, []string, error) {
	// Implement the parsing logic for Socket.IO messages
	var socketIOMessage SocketIOMessage
	if err := json.Unmarshal(message, &socketIOMessage); err != nil {
		return "", nil, err
	}

	messageType := strconv.Itoa(socketIOMessage.Type)
	data := socketIOMessage.Data

	return messageType, data, nil
}

func sendSocketIOMessage(socket *websocket.Conn, event string, data interface{}) error {
	// Create a Socket.IO message and send it over the WebSocket connection
	message := SocketIOMessage{
		Type: 2, // Message type for event
		Data: []string{event, toJSONString(data)},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return socket.WriteMessage(websocket.TextMessage, messageBytes)
}

func toJSONString(data interface{}) string {
	// Convert data to a JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return ""
	}
	return string(jsonData)
}

func handleCustomEvent(data []string) {
	// Handle custom events received from the Socket.IO server
	if len(data) < 2 {
		log.Println("Received invalid custom event data:", data)
		return
	}

	event := data[0]
	eventData := data[1]

	fmt.Printf("Received custom event '%s' with data: %s\n", event, eventData)
}
