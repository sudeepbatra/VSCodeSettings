package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketURL = "ws://smartapisocket.angelone.in/smart-stream"
)

type Request struct {
	CorrelationID string        `json:"correlationID,omitempty"`
	Action        int           `json:"action"`
	Params        RequestParams `json:"params"`
}

type RequestParams struct {
	Mode       int                `json:"mode"`
	TokenLists []RequestTokenList `json:"tokenList"`
}

type RequestTokenList struct {
	ExchangeType int      `json:"exchangeType"`
	Tokens       []string `json:"tokens"`
}

func main() {
	dialer := websocket.Dialer{}
	headers := http.Header{
		"Authorization": {"Bearer eyJhbGciOiJIUzUxMiJ9.eyJ0b2tlbiI6IlJFRlJFU0gtVE9LRU4iLCJpYXQiOjE2OTI1MDU2MDJ9.AzQHeoHxgjlUrAIGX2HUfIzTZnCkMR6mgH5fxy_mULG8IMdUHDJddjEEjgC2Om8_9AtAs9fvf4tqLKwpKKpz2g"},
		"x-api-key":     {"iGKWS2zU"},
		"x-client-code": {"S1632585"},
		"x-feed-token":  {"eyJhbGciOiJIUzUxMiJ9.eyJ1c2VybmFtZSI6IlMxNjMyNTg1IiwiaWF0IjoxNjkyNTA1NTk4LCJleHAiOjE2OTI1OTE5OTh9.KGZgCcD6w1v6FYW1FraRUu4Ngfw8ffXAVfiYb-fR9G64XvJ8mbSgl6zhADolO_V7NerxNUfeSf_TeN4u9yg4TA"},
	}
	conn, _, err := dialer.Dial(websocketURL, headers)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	err = conn.WriteJSON(authHeaders)
	if err != nil {
		log.Println("Error sending auth headers:", err)
		return
	}

	// Heartbeat message
	go sendHeartbeats(conn)

	// Subscribe to tokens
	subscribe(conn)
}

func sendHeartbeats(conn *websocket.Conn) {
	for {
		time.Sleep(30 * time.Second)
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			log.Println("Error sending heartbeat:", err)
			return
		}
	}
}

func subscribe(conn *websocket.Conn) {
	request := Request{
		CorrelationID: "abcde12345",
		Action:        1, // Subscribe action
		Params: RequestParams{
			Mode: 1, // Subscription mode (LTP)
			TokenLists: []RequestTokenList{
				{
					ExchangeType: 1, // nse_cm
					Tokens:       []string{"10626", "5290"},
				},
				{
					ExchangeType: 5, // mcx_fo
					Tokens:       []string{"234230", "234235", "234219"},
				},
			},
		},
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshaling request:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		log.Println("Error sending subscribe request:", err)
		return
	}

	// Handle response
	for {
		_, respBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading response:", err)
			return
		}

		// Process the binary response
		processResponse(respBytes)
	}
}

func processResponse(respBytes []byte) {
	// Parse the binary response based on the given specifications
	// Implement the actual parsing logic here using the binary package
	// Extract and interpret the data fields based on their types and positions
	// Handle different response types (LTP, Quote, Snap Quote) accordingly
	fmt.Println("Received response:", respBytes)
}

func byteArrayToLittleEndianInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}

func byteArrayToLittleEndianInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

// Add more conversion functions as needed
