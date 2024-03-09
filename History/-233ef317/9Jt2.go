package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	socketio "github.com/ambelovsky/gosf-socketio"
)

type SocketManager struct {
	socket       *socketio.Client
	socketOrder  *socketio.Client
	socketOHLC   *socketio.Client
	userId       string
	sessionToken string
	apiKey       string
	config       map[string]string // Assuming you have a config map
}

func (sm *SocketManager) connect(isOrder bool, isOHLC bool) {
	authDict := map[string]string{
		"user":   sm.userId,
		"token":  sm.sessionToken,
		"appkey": sm.apiKey,
	}

	if isOrder && sm.socketOrder == nil {
		options := &socketio.ClientOptions{
			Transports: []string{"websocket"},
			Query:      url.Values(authDict).Encode(),
		}
		uriOrder := sm.config["LIVE_FEEDS_URL"]
		sm.socketOrder, _ = socketio.NewClient(strings.TrimSuffix(uriOrder, "/")+"/socket.io/", options)
		sm.socketOrder.On("connect", func() {
			fmt.Println("Connected to Order Socket.IO")
		})
		sm.socketOrder.Connect()
	}

	if isOHLC && sm.socketOHLC == nil {
		options := &socketio.ClientOptions{
			Transports: []string{"websocket"},
			Query:      url.Values(authDict).Encode(),
			Path:       "/ohlcvstream/",
		}
		uriOHLC := sm.config["LIVE_OHLC_STREAM_URL"]
		sm.socketOHLC, _ = socketio.NewClient(strings.TrimSuffix(uriOHLC, "/")+"/socket.io/", options)
		sm.socketOHLC.On("connect", func() {
			fmt.Println("Connected to OHLC Socket.IO")
		})
		sm.socketOHLC.Connect()
	}

	if sm.socket == nil {
		options := &socketio.ClientOptions{
			Transports: []string{"websocket"},
			Query:      url.Values(authDict).Encode(),
			Headers:    http.Header{"User-Agent": []string{"node-socketio[client]/socket"}},
		}
		uri := sm.config["LIVE_STREAM_URL"]
		sm.socket, _ = socketio.NewClient(strings.TrimSuffix(uri, "/")+"/socket.io/", options)
		sm.socket.On("connect", func() {
			fmt.Println("Connected to Socket.IO")
		})
		sm.socket.Connect()
	}
}

func (sm *SocketManager) connectTicker() {
	sm.connect(false, false)
}

func (sm *SocketManager) watch(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("joining stock subscription:", stock)
			sm.socket.Emit("join", stock)
		}
	}
}

func (sm *SocketManager) unWatch(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("leaving stock subscription:", stock)
			sm.socket.Emit("leave", stock)
		}
	}
}

func (sm *SocketManager) watchOHLC(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("joining OHLC subscription:", stock)
			sm.socketOHLC.Emit("join", stock)
		}
	}
}

func (sm *SocketManager) unWatchOHLC(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("leaving OHLC subscription:", stock)
			sm.socketOHLC.Emit("leave", stock)
		}
	}
}

func (sm *SocketManager) subscribeFeeds(stockToken string) error {
	if stockToken == "" {
		return fmt.Errorf("Blank stock token")
	}
	sm.watch([]string{stockToken})
	return nil
}

func (sm *SocketManager) subscribeFeedsWithInterval(stockToken, interval string) error {
	if interval == "" || !contains(sm.config["INTERVAL_TYPES_STREAM_OHLC"], interval) {
		return fmt.Errorf("Invalid interval type")
	}
	// sm.interval = sm.config["channelIntervalMap"][interval]
	sm.connect(false, true)
	sm.watchOHLC([]string{stockToken})
	return nil
}

func contains(list string, item string) bool {
	return strings.Contains(list, item)
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
	err := sm.subscribeFeeds("your_stock_token")
	if err != nil {
		fmt.Println("Error subscribing to feeds:", err)
	}

	err = sm.subscribeFeedsWithInterval("your_stock_token", "your_interval")
	if err != nil {
		fmt.Println("Error subscribing to feeds with interval:", err)
	}

	// Keep the program running
	select {}
}
