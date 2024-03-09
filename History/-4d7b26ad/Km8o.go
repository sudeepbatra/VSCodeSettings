package fivepaisa

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	accessToken              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6IjUwNjAzNzEwIiwicm9sZSI6IjEwNDI3IiwiU3RhdGUiOiIiLCJSZWRpcmVjdFNlcnZlciI6IkEiLCJuYmYiOjE2OTMyODA1MDIsImV4cCI6MTY5MzMzMzc5OSwiaWF0IjoxNjkzMjgwNTAyfQ.NsJO9Xu5hY0y56OSFex8nUM6R_LqYh7d0W6cQGRbu5k"
	MarketFeedMethod         = "MarketFeedV3"
	MarketDepthServiceMethod = "MarketDepthService"
	GetOpenInterestMethod    = "GetScripInfoForFuture"
	IndicesMethod            = "Indices"
	SubscribeOperation       = "Subscribe"
	UnsubscribeOperation     = "Unsubscribe"
)

type WebsocketFivePaisa struct {
	*websocket.Conn
	clientCode  string
	accessToken string
}

func (wfp *WebsocketFivePaisa) IsConnected() bool {
	return wfp.Conn != nil
}

func (wfp *WebsocketFivePaisa) Close() error {
	err := wfp.Conn.Close()
	wfp.Conn = nil

	return err
}

func (wfp WebsocketFivePaisa) Reconnect(url string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	wfp.Conn = conn

	return nil
}

func connectWebSocket(socketURL string) (*WebsocketFivePaisa, error) {
	wfp := &WebsocketFivePaisa{}
	err := wfp.Reconnect(socketURL)
	if err != nil {
		return nil, err
	}

	return wfp, nil
}

// func connectWebSocket() (*websocket.Conn, error) {
// 	socketURL := fmt.Sprintf("wss://openfeed.5paisa.com/Feeds/api/chat?Value1=%s|%s",
// 		url.QueryEscape(accessToken), url.QueryEscape(clientCode))

// 	var dialer websocket.Dialer
// 	dialer.Subprotocols = []string{"chat"}

// 	conn, _, err := dialer.Dial(socketURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

type Subscription struct {
	Exch      string `json:"Exch"`
	ExchType  string `json:"ExchType"`
	ScripCode int    `json:"ScripCode"`
}

type WSMessage struct {
	Method         string         `json:"Method"`
	Operation      string         `json:"Operation"`
	ClientCode     string         `json:"ClientCode"`
	MarketFeedData []Subscription `json:"MarketFeedData"`
}

func readMessages(wfp WebsocketFivePaisa) {
	for wfp.IsConnected() {
		_, message, err := wfp.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			wfp.Close()
			return
		}

		log.Info("Received:", string(message))
	}
}

// func readMessages(conn *websocket.Conn) {
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("WebSocket read error:", err)
// 			return
// 		}

// 		log.Info("Received:", string(message))
// 	}
// }

func WebsocketConnect() {
	socketURL := fmt.Sprintf("wss://openfeed.5paisa.com/Feeds/api/chat?Value1=%s|%s",
		url.QueryEscape(accessToken), url.QueryEscape(clientCode))

	conn, err := connectWebSocket(socketURL)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
		return
	}
	defer conn.Close()

	wsPayload := sampleSubscription()

	err = conn.WriteJSON(wsPayload)
	if err != nil {
		log.Println("WebSocket write error:", err)
		return
	}

	go readMessages(wfp)

	for {
		time.Sleep(time.Second)
	}
}

func sampleSubscription() WSMessage {
	sampleSubscriptions := []Subscription{
		{Exch: "N", ExchType: "C", ScripCode: 15083},
		{Exch: "B", ExchType: "C", ScripCode: 999901},
		{Exch: "N", ExchType: "C", ScripCode: 22},
	}

	wsPayload := WSMessage{
		Method:         MarketFeedMethod,
		Operation:      SubscribeOperation,
		ClientCode:     clientCode,
		MarketFeedData: sampleSubscriptions,
	}

	return wsPayload
}
