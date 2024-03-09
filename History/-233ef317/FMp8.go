package main

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
)

type SocketManager struct {
	UserID       string
	SessionToken string
	APIKey       string
	Socket       socketio.Server
	SocketOrder  socketio.Server
	SocketOHLC   socketio.Server
}

func NewSocketManager(userID, sessionToken, apiKey string) *SocketManager {
	return &SocketManager{
		UserID:       userID,
		SessionToken: sessionToken,
		APIKey:       apiKey,
	}
}

func (sm *SocketManager) Connect(isOrder, isOHLC bool) {
	authDict := map[string]interface{}{
		"user":   sm.UserID,
		"token":  sm.SessionToken,
		"appkey": sm.APIKey,
	}

	if isOrder && sm.SocketOrder == nil {
		server, err := sm.createSocketServer("/socket.io/order")
		if err != nil {
			fmt.Println("Error creating order socket:", err)
			return
		}
		sm.SocketOrder = server
	}

	if isOHLC && sm.SocketOHLC == nil {
		server, err := sm.createSocketServer("/socket.io/ohlcvstream")
		if err != nil {
			fmt.Println("Error creating OHLC socket:", err)
			return
		}
		sm.SocketOHLC = server
	} else if sm.Socket == nil {
		server, err := sm.createSocketServer("/socket.io/live")
		if err != nil {
			fmt.Println("Error creating live socket:", err)
			return
		}
		sm.Socket = server
	}
}

func (sm *SocketManager) createSocketServer(path string) (socketio.Server, error) {
	server, err := socketio.NewServer(&engineio.Options{
		Transports: []string{"websocket"},
	})
	if err != nil {
		return nil, err
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		// Handle connection logic here
		s.SetContext(authDict)
		return nil
	})

	// Add event handlers here
	// Example:
	// server.OnEvent("/", "join", func(s socketio.Conn, stock string) {
	//     // Handle "join" event
	// })

	go server.Serve()
	return server, nil
}

func (sm *SocketManager) Watch(stocks []string) {
	for _, stock := range stocks {
		fmt.Println("joining stock subscription:", stock)
		sm.Socket.Emit("join", stock)
	}
}

func (sm *SocketManager) UnWatch(stocks []string) {
	for _, stock := range stocks {
		fmt.Println("leaving stock subscription:", stock)
		sm.Socket.Emit("leave", stock)
	}
}

func (sm *SocketManager) WatchOHLC(stocks []string) {
	for _, stock := range stocks {
		sm.SocketOHLC.Emit("join", stock)
	}
}

func (sm *SocketManager) UnWatchOHLC(stocks []string) {
	for _, stock := range stocks {
		sm.SocketOHLC.Emit("leave", stock)
	}
}

func (sm *SocketManager) SubscribeFeeds(stockToken string) {
	if stockToken == "" {
		fmt.Println("Stock token is blank")
		return
	}

	sm.Watch([]string{stockToken})
	fmt.Printf("Subscribed to stock feeds for: %s\n", stockToken)
}

func (sm *SocketManager) SubscribeFeedsWithInterval(stockToken, interval string) {
	if interval != "" {
		// Check if the interval is valid, implement the interval validation logic here
		// Example:
		// if !isValidInterval(interval) {
		//     fmt.Println("Invalid interval")
		//     return
		// }
	} else {
		fmt.Println("Interval is blank")
		return
	}

	sm.Connect(false, true)

	if stockToken == "" {
		fmt.Println("Stock token is blank")
		return
	}

	sm.WatchOHLC([]string{stockToken})
	fmt.Printf("Subscribed to OHLC feeds for: %s\n", stockToken)
}

func main() {
	userID := "your_user_id"
	sessionToken := "your_session_token"
	apiKey := "your_api_key"

	sm := NewSocketManager(userID, sessionToken, apiKey)
	sm.ConnectTicker()

	// Example usage:
	sm.SubscribeFeeds("AAPL")
	sm.SubscribeFeedsWithInterval("AAPL", "1min")
}
