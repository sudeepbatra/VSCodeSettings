package main

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

type MySocketClient struct {
	userId       string
	sessionToken string
	apiKey       string
	socket       *socketio.Server
	socketOrder  *socketio.Server
	socketOHLC   *socketio.Server
}

func NewMySocketClient(userId, sessionToken, apiKey string) *MySocketClient {
	return &MySocketClient{
		userId:       userId,
		sessionToken: sessionToken,
		apiKey:       apiKey,
	}
}

func (c *MySocketClient) connect(isOrder, isOHLC bool) {
	authDict := map[string]string{
		"user":   c.userId,
		"token":  c.sessionToken,
		"appkey": c.apiKey,
	}

	if isOrder && c.socketOrder == nil {
		c.socketOrder = createSocketIOClient("https://livefeeds.icicidirect.com", authDict)
	}

	if isOHLC && c.socketOHLC == nil {
		c.socketOHLC = createSocketIOClient("https://breezeapi.icicidirect.com", authDict)
	} else if c.socket == nil {
		c.socket = createSocketIOClient("https://livestream.icicidirect.com", authDict)
	}
}

func createSocketIOClient(serverURL string, authDict map[string]string) *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("Socket.IO connection error:", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		// Authenticate here using authDict
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Handle disconnect
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal("Socket.IO server error:", err)
		}
	}()

	return server
}

func (c *MySocketClient) watch(stocks []string, server *socketio.Server) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("joining stock subscription:", stock)
			server.BroadcastTo("/", "join", stock)
		}
	}
}

func (c *MySocketClient) unWatch(stocks []string, server *socketio.Server) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("leaving stock subscription:", stock)
			server.BroadcastTo("/", "leave", stock)
		}
	}
}

func (c *MySocketClient) subscribeFeeds(stockToken string, server *socketio.Server) {
	if stockToken == "" {
		return
	}

	c.watch([]string{stockToken}, server)
	// Handle response as needed
}

func (c *MySocketClient) subscribeFeedsWithInterval(stockToken, interval string, server *socketio.Server) {
	if interval == "" {
		return
	}

	// Handle interval and subscription as needed
}

func main() {
	// Replace with your actual credentials
	userId := "your_user_id"
	sessionToken := "your_session_token"
	apiKey := "your_api_key"

	client := NewMySocketClient(userId, sessionToken, apiKey)
	client.connect(true, true) // Example: Connect to both order and OHLC streams

	// Subscribe to feeds
	client.subscribeFeeds("your_stock_token", client.socket)
	// client.subscribeFeedsWithInterval("your_stock_token", "your_interval", client.socket)

	// Keep the program running
	select {}
}
