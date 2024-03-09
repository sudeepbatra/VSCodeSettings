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
	// if len(respBytes) < 379 {
	// 	log.Println("Received incomplete response.")
	// 	return
	// }

	parsedData := make(map[string]interface{})
	parsedData["subscription_mode"] = unpackData(binaryData, 0, 1, "B")[0]
	parsedData["exchange_type"] = unpackData(binaryData, 1, 2, "B")[0]
	parsedData["token"] = parseTokenValue(binaryData[2:27])
	parsedData["sequence_number"] = unpackData(binaryData, 27, 35, "q")[0]
	parsedData["exchange_timestamp"] = unpackData(binaryData, 35, 43, "q")[0]
	parsedData["last_traded_price"] = unpackData(binaryData, 43, 51, "q")[0]

	subscriptionModeVal, found := subscriptionModeMap[parsedData["subscription_mode"].(int)]
	if found {
		parsedData["subscription_mode_val"] = subscriptionModeVal
	}

	subMode := parsedData["subscription_mode"].(int)
	if subMode == quote || subMode == snapQuote {
		parsedData["last_traded_quantity"] = unpackData(binaryData, 51, 59, "q")[0]
		parsedData["average_traded_price"] = unpackData(binaryData, 59, 67, "q")[0]
		parsedData["volume_trade_for_the_day"] = unpackData(binaryData, 67, 75, "q")[0]
		parsedData["total_buy_quantity"] = unpackData(binaryData, 75, 83, "d")[0]
		parsedData["total_sell_quantity"] = unpackData(binaryData, 83, 91, "d")[0]
		parsedData["open_price_of_the_day"] = unpackData(binaryData, 91, 99, "q")[0]
		parsedData["high_price_of_the_day"] = unpackData(binaryData, 99, 107, "q")[0]
		parsedData["low_price_of_the_day"] = unpackData(binaryData, 107, 115, "q")[0]
		parsedData["closed_price"] = unpackData(binaryData, 115, 123, "q")[0]
	}

	if subMode == snapQuote {
		parsedData["last_traded_timestamp"] = unpackData(binaryData, 123, 131, "q")[0]
		parsedData["open_interest"] = unpackData(binaryData, 131, 139, "q")[0]
		parsedData["open_interest_change_percentage"] = unpackData(binaryData, 139, 147, "q")[0]
		parsedData["upper_circuit_limit"] = unpackData(binaryData, 347, 355, "q")[0]
		parsedData["lower_circuit_limit"] = unpackData(binaryData, 355, 363, "q")[0]
		parsedData["52_week_high_price"] = unpackData(binaryData, 363, 371, "q")[0]
		parsedData["52_week_low_price"] = unpackData(binaryData, 371, 379, "q")[0]
		best5BuyAndSellData := parseBest5BuyAndSellData(binaryData[147:347])
		parsedData["best_5_buy_data"] = best5BuyAndSellData["best_5_sell_data"]
		parsedData["best_5_sell_data"] = best5BuyAndSellData["best_5_buy_data"]
	}

}

func byteArrayToLittleEndianInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}

func byteArrayToLittleEndianInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

func unpackData(binaryData []byte, start int, end int, byteFormat string) []interface{} {
	dataSlice := binaryData[start:end]
	var result []interface{}

	switch byteFormat {
	case "B":
		result = append(result, uint8(dataSlice[0]))
	case "q":
		result = append(result, int64(binary.LittleEndian.Uint64(dataSlice)))
	case "d":
		result = append(result, float64(binary.LittleEndian.Uint64(dataSlice)))
		// Add more cases for other byte formats as needed
	}

	return result
}

// Add more conversion functions as needed
