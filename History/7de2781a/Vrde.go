package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Angel One Smart API Websocket Streaming 2.0
	rootURI               = "ws://smartapisocket.angelone.in/smart-stream"
	heartBeatMessage      = "ping"
	heartBeatInterval     = 10 * time.Second
	littleEndianByteOrder = "<"

	// Available Actions
	SubscribeAction   = 1
	UnsubscribeAction = 0

	// Possible Subscribtion Modes
	LTPMode   = 1
	Quote     = 2
	SnapQuote = 3

	//Exchange Types
	NseCm = 1
	NseFo = 2
	BseCm = 3
	BseFo = 4
	McxFo = 5
	NcxFo = 7
	CdeFo = 13
)

var (
	SUBSCRIPTION_MODE_MAP = map[int]string{
		1: "LTP",
		2: "QUOTE",
		3: "SNAP_QUOTE",
	}
	resubscribeFlag = false
	upgrader        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Subscription struct {
	CorrelationID string        `json:"correlationID"`
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

type ErrorResponse struct {
	CorrelationID string `json:"correlationID"`
	ErrorCode     string `json:"errorCode"`
	ErrorMessage  string `json:"errorMessage"`
}

type ParsedData struct {
	SubscriptionMode             uint8
	ExchangeType                 uint8
	Token                        string
	SequenceNumber               int64
	ExchangeTimestamp            int64
	LastTradedPrice              int64
	SubscriptionModeVal          string
	LastTradedQuantity           int64
	AverageTradedPrice           int64
	VolumeTradeForTheDay         int64
	TotalBuyQuantity             float64
	TotalSellQuantity            float64
	OpenPriceOfTheDay            int64
	HighPriceOfTheDay            int64
	LowPriceOfTheDay             int64
	ClosedPrice                  int64
	LastTradedTimestamp          int64
	OpenInterest                 int64
	OpenInterestChangePercentage int64
	UpperCircuitLimit            int64
	LowerCircuitLimit            int64
	Week52HighPrice              int64
	Week52LowPrice               int64
	Best5BuyData, Best5SellData  []int64 // Modify the data types as required
}

type SmartWebSocketV2 struct {
	wsapp               *websocket.Conn
	authToken           string
	apiKey              string
	clientCode          string
	feedToken           string
	maxRetryAttempt     int
	disconnectFlag      bool
	lastPongTimestamp   int64
	inputRequestDict    map[int]map[int][]string // Mode -> ExchangeType -> Tokens
	currentRetryAttempt int
}

func NewSmartWebSocketV2(authToken, apiKey, clientCode, feedToken string, maxRetryAttempt int) *SmartWebSocketV2 {
	return &SmartWebSocketV2{
		authToken:        authToken,
		apiKey:           apiKey,
		clientCode:       clientCode,
		feedToken:        feedToken,
		maxRetryAttempt:  maxRetryAttempt,
		disconnectFlag:   true,
		inputRequestDict: make(map[int]map[int][]string),
	}
}

func main() {
	headers := http.Header{
		"Authorization": {"Bearer eyJhbGciOiJIUzUxMiJ9.eyJ0b2tlbiI6IlJFRlJFU0gtVE9LRU4iLCJpYXQiOjE2OTI1MDU2MDJ9.AzQHeoHxgjlUrAIGX2HUfIzTZnCkMR6mgH5fxy_mULG8IMdUHDJddjEEjgC2Om8_9AtAs9fvf4tqLKwpKKpz2g"},
		"x-api-key":     {"iGKWS2zU"},
		"x-client-code": {"S1632585"},
		"x-feed-token":  {"eyJhbGciOiJIUzUxMiJ9.eyJ0b2tlbiI6IlJFRlJFU0gtVE9LRU4iLCJpYXQiOjE2OTI1OTMxMjh9.hw7hdWM1ez0qoqCY4GgiHAYQwM6nrrJQqQchwNIHfY2EPHOXnJvycCFXbB0JR5ZZLGAQpBMWoGdC9gP_Ea6j3Q"
		FeedToken:},
	}
	conn, _, err := websocket.DefaultDialer.Dial(rootURI, headers)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	// Heartbeat message
	go sendHeartbeats(conn)

	subscribe(conn)
}

func subscribe(conn *websocket.Conn) {
	request := Subscription{
		CorrelationID: "abcde12345",
		Action:        SubscribeAction, // Subscribe action
		Params: RequestParams{
			Mode: SnapQuote, // Subscription mode (LTP)
			TokenLists: []RequestTokenList{
				{
					ExchangeType: NseCm, // nse_cm
					Tokens:       []string{"2885", "5290"},
				},
				{
					ExchangeType: McxFo, // mcx_fo
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

func sendHeartbeats(conn *websocket.Conn) {
	for {
		time.Sleep(heartBeatInterval)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(heartBeatMessage)); err != nil {
			fmt.Println("Error sending heartbeat:", err)
			return
		}
	}
}

// func handleTextMessage(conn *websocket.Conn, message []byte) {
// 	// Parse the received message and handle accordingly
// 	var subscription Subscription
// 	err := json.Unmarshal(message, &subscription)
// 	if err != nil {
// 		fmt.Println("Error parsing subscription request:", err)
// 		return
// 	}

// 	// Process the subscription request
// 	// ...

// 	// Send the subscription response
// 	response := []byte(`{"correlationID": "abcde12345", "response": "Subscription successful"}`)
// 	if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
// 		fmt.Println("Error sending response:", err)
// 		return
// 	}
// }

func processResponse(respBytes []byte) {
	fmt.Println("Received response:", respBytes)

	var parsedData ParsedData

	if bytes.Equal(respBytes, []byte("pong")) {
		fmt.Println("respBytes is 'pong'")
		return
	}

	// Implement _unpack_data function in Go to unpack data from binaryData using binary.Read

	parsedData.SubscriptionMode = respBytes[0]
	parsedData.ExchangeType = respBytes[1]
	parsedData.Token = parseTokenValue(respBytes[2:27])
	parsedData.SequenceNumber, _ = unpackData(respBytes[27:35])
	parsedData.ExchangeTimestamp, _ = unpackData(respBytes[35:43])
	parsedData.LastTradedPrice, _ = unpackData(respBytes[43:51])
	fmt.Println("1 ParsedData: ", parsedData)

	if parsedData.SubscriptionMode == Quote || parsedData.SubscriptionMode == SnapQuote {
		parsedData.LastTradedQuantity, _ = unpackData(respBytes[51:59])
		parsedData.AverageTradedPrice, _ = unpackData(respBytes[59:67])
		parsedData.VolumeTradeForTheDay, _ = unpackData(respBytes[67:75])
		parsedData.TotalBuyQuantity, _ = unpackDoubleData(respBytes[75:83])
		parsedData.TotalSellQuantity, _ = unpackDoubleData(respBytes[83:91])
		parsedData.OpenPriceOfTheDay, _ = unpackData(respBytes[91:99])
		parsedData.HighPriceOfTheDay, _ = unpackData(respBytes[99:107])
		parsedData.LowPriceOfTheDay, _ = unpackData(respBytes[107:115])
		parsedData.ClosedPrice, _ = unpackData(respBytes[115:123])
	}

	if parsedData.SubscriptionMode == SnapQuote {
		parsedData.LastTradedTimestamp, _ = unpackData(respBytes[123:131])
		parsedData.OpenInterest, _ = unpackData(respBytes[131:139])
		parsedData.OpenInterestChangePercentage, _ = unpackData(respBytes[139:147])
		parsedData.UpperCircuitLimit, _ = unpackData(respBytes[347:355])
		parsedData.LowerCircuitLimit, _ = unpackData(respBytes[355:363])
		parsedData.Week52HighPrice, _ = unpackData(respBytes[363:371])
		parsedData.Week52LowPrice, _ = unpackData(respBytes[371:379])

		// Implement _parseBest5BuyAndSellData function in Go

		// Modify best_5_buy_and_sell_data to return two slices of int64 data
		// best5BuyData, best5SellData := parseBest5BuyAndSellData(binaryData[147:347])
		// parsedData.Best5BuyData = best5BuyData
		// parsedData.Best5SellData = best5SellData
	}

	fmt.Println("2 ParsedData: ", parsedData)
}

func parseTokenValue(binaryPacket []byte) string {
	var token bytes.Buffer

	for _, b := range binaryPacket {
		if b == 0 {
			break
		}
		token.WriteByte(b)
	}

	return token.String()
}

func unpackData(binaryData []byte) (value int64, err error) {
	err = binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("---Error---", err)
		return
	}

	return
}

func unpackDoubleData(binaryData []byte) (value float64, err error) {
	err = binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("---Error---", err)
		return
	}

	return
}

// type Best5Data struct {
// 	Flag       uint16
// 	Quantity   int64
// 	Price      int64
// 	NoOfOrders uint16
// }

// func ParseBest5BuyAndSellData(binaryData []byte) (map[string][]Best5Data, error) {
// 	best5BuyAndSellData := p.parseBest5BuyAndSellData(binaryData)
// 	parsedData := make(map[string][]Best5Data)
// 	parsedData["best_5_buy_data"] = best5BuyAndSellData["best_5_sell_data"]
// 	parsedData["best_5_sell_data"] = best5BuyAndSellData["best_5_buy_data"]
// 	return parsedData, nil
// }

// func parseBest5BuyAndSellData(binaryData []byte) map[string][]Best5Data {
// 	buySellPackets := p.splitPackets(binaryData)
// 	best5BuyData := make([]Best5Data, 0)
// 	best5SellData := make([]Best5Data, 0)

// 	for _, packet := range buySellPackets {
// 		eachData := Best5Data{
// 			Flag:       binary.BigEndian.Uint16(packet[0:2]),
// 			Quantity:   int64(binary.BigEndian.Uint64(packet[2:10])),
// 			Price:      int64(binary.BigEndian.Uint64(packet[10:18])),
// 			NoOfOrders: binary.BigEndian.Uint16(packet[18:20]),
// 		}

// 		if eachData.Flag == 0 {
// 			best5BuyData = append(best5BuyData, eachData)
// 		} else {
// 			best5SellData = append(best5SellData, eachData)
// 		}
// 	}

// 	return map[string][]Best5Data{
// 		"best_5_buy_data":  best5BuyData,
// 		"best_5_sell_data": best5SellData,
// 	}
// }

// func splitPackets(binaryPackets []byte) [][]byte {
// 	packetSize := 20
// 	numPackets := len(binaryPackets) / packetSize
// 	packets := make([][]byte, numPackets)

// 	for i := 0; i < numPackets; i++ {
// 		start := i * packetSize
// 		end := start + packetSize
// 		packets[i] = binaryPackets[start:end]
// 	}

// 	return packets
// }
