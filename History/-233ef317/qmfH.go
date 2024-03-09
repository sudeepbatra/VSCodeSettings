package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	socketio "github.com/zhouhui8915/go-socket.io-client"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SocketManager struct {
	socket       *socketio.Client
	socketOrder  *socketio.Client
	socketOHLC   *socketio.Client
	userId       string
	sessionToken string
	apiKey       string
	config       map[string]string // Assuming you have a config map
}

func (sm *SocketManager) sendJoin(c *socketio.Client) {
	log.Println("Emitting /join")
	err := c.Emit("/join", Channel{"main"})
	if err != nil {
		log.Fatal(err)
	}
}

func (sm *SocketManager) connect(isOrder, isOHLC bool) {
	authDict := map[string]string{
		"user":   sm.userId,
		"token":  sm.sessionToken,
		"appkey": sm.apiKey,
	}

	if isOrder && sm.socketOrder == nil {
		options := &socketio.Options{
			Transport: "websocket",
			Query:     url.Values(authDict).Encode(),
		}
		uriOrder := sm.config["LIVE_FEEDS_URL"]
		sm.socketOrder, _ = socketio.NewClient(strings.TrimSuffix(uriOrder, "/")+"/socket.io/", options)
		sm.socketOrder.On("connect", func() {
			fmt.Println("Connected to Order Socket.IO")
		})
		sm.socketOrder.Connect()
	}

	if isOHLC && sm.socketOHLC == nil {
		options := &socketio.Options{
			Transport: "websocket",
			Query:     url.Values(authDict).Encode(),
			Path:      "/ohlcvstream/",
		}
		uriOHLC := sm.config["LIVE_OHLC_STREAM_URL"]
		sm.socketOHLC, _ = socketio.NewClient(strings.TrimSuffix(uriOHLC, "/")+"/socket.io/", options)
		sm.socketOHLC.On("connect", func() {
			fmt.Println("Connected to OHLC Socket.IO")
		})
		sm.socketOHLC.Connect()
	}

	if sm.socket == nil {
		options := &socketio.Options{
			Transport: "websocket",
			Query:     url.Values(authDict).Encode(),
			Headers:   http.Header{"User-Agent": []string{"node-socketio[client]/socket"}},
		}
		uri := sm.config["LIVE_STREAM_URL"]
		sm.socket, _ = socketio.NewClient(strings.TrimSuffix(uri, "/")+"/socket.io/", options)
		sm.socket.On("connect", func() {
			fmt.Println("Connected to Socket.IO")
		})
		sm.socket.Connect()
	}
}

func main() {
	sm := &SocketManager{
		userId:       "your_user_id",
		sessionToken: "your_session_token",
		apiKey:       "your_api_key",
		config: map[string]string{
			"LIVE_FEEDS_URL":             "https://livefeeds.icicidirect.com",
			"LIVE_STREAM_URL":            "https://livestream.icicidirect.com",
			"LIVE_OHLC_STREAM_URL":       "https://breezeapi.icicidirect.com",
			"INTERVAL_TYPES_STREAM_OHLC": "your_interval_types",
			"channelIntervalMap":         "your_channel_interval_map",
		},
	}

	sm.connectTicker()

	// Keep the program running
	select {}
}
