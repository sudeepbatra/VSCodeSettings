package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Socket represents a Socket.IO client in Go.
type Socket struct {
	conn *websocket.Conn
}

// NewSocket creates a new Socket.IO client.
func NewSocket(serverURL string) (*Socket, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &Socket{conn: conn}, nil
}

// Connect establishes a connection to the server.
func (s *Socket) Connect() {
	// Implement your connection logic here
}

// Send sends a message to the server.
func (s *Socket) Send(message string) {
	// Implement your message sending logic here
}

// Close closes the WebSocket connection.
func (s *Socket) Close() {
	err := s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error sending close message:", err)
	}

	err = s.conn.Close()
	if err != nil {
		log.Println("Error closing WebSocket connection:", err)
	}
}

func main() {
	serverURL := "ws://your-socket-io-server-url"
	socket, err := NewSocket(serverURL)
	if err != nil {
		log.Fatal("Error creating Socket.IO client:", err)
	}

	// Handle Ctrl+C to gracefully close the WebSocket connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			select {
			case <-interrupt:
				fmt.Println("Received interrupt signal. Closing WebSocket...")
				socket.Close()
				return
			}
		}
	}()

	// Start the Socket.IO connection
	socket.Connect()

	// Send messages to the server
	for {
		message := "Hello, Socket.IO Server!"
		socket.Send(message)
		time.Sleep(5 * time.Second)
	}
}
