package main

import (
	"fmt"
	"log"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

type MySocketClient struct {
	userId       string
	sessionToken string
	apiKey       string
	socket       *socketio_client.Client
	socketOrder  *socketio_client.Client
	socketOHLC   *socketio_client.Client
}

func NewMySocketClient(userId, sessionToken, apiKey string) *MySocketClient {
	return &MySocketClient{
		userId:       userId,
		sessionToken: sessionToken,
		apiKey:       apiKey,
	}
}

func (c *MySocketClient) connect(isOrder, isOHLC bool) {
	authDict := map[string]interface{}{
		"user":   c.userId,
		"token":  c.sessionToken,
		"appkey": c.apiKey,
	}

	uriOrder := "https://livefeeds.icicidirect.com" // Replace with your order feed URL
	uriOHLC := "https://breezeapi.icicidirect.com"  // Replace with your OHLC feed URL
	uri := "https://livestream.icicidirect.com"     // Replace with your default feed URL

	var err error

	if isOrder && c.socketOrder == nil {
		c.socketOrder, err = createSocketIOClient(uriOrder, authDict)
		if err != nil {
			log.Fatal("Socket.IO connection error:", err)
		}
	}

	if isOHLC && c.socketOHLC == nil {
		c.socketOHLC, err = createSocketIOClient(uriOHLC, authDict)
		if err != nil {
			log.Fatal("Socket.IO connection error:", err)
		}
	} else if c.socket == nil {
		c.socket, err = createSocketIOClient(uri, authDict)
		if err != nil {
			log.Fatal("Socket.IO connection error:", err)
		}
	}
}

func createSocketIOClient(uri string, authDict map[string]interface{}) (*socketio_client.Client, error) {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	// Convert authDict values to strings
	for key, value := range authDict {
		opts.Query[key] = fmt.Sprintf("%v", value)
	}

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		return nil, err
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connect", func() {
		log.Printf("on connect\n")
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	return client, nil
}

func (c *MySocketClient) watch(stocks []string, client *socketio_client.Client) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("joining stock subscription:", stock)
			client.Emit("join", stock)
		}
	}
}

func (c *MySocketClient) unWatch(stocks []string, client *socketio_client.Client) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("leaving stock subscription:", stock)
			client.Emit("leave", stock)
		}
	}
}

func (c *MySocketClient) subscribeFeeds(stockToken string) {
	if stockToken == "" {
		return
	}

	c.watch([]string{stockToken}, c.socket)
	// Handle response as needed
}

func (c *MySocketClient) subscribeFeedsWithInterval(stockToken, interval string) {
	if interval == "" {
		return
	}

	// Handle interval and subscription as needed
}

func main() {
	// Replace with your actual credentials
	userId := "SUDDEPBA"
	sessionToken := "20340551"
	apiKey := "s76162#+U35414Y*S413=099_FA6P567"

	client := NewMySocketClient(userId, sessionToken, apiKey)
	client.connect(true, true) // Example: Connect to both order and OHLC streams

	// Subscribe to feeds
	client.subscribeFeeds("your_stock_token")
	// client.subscribeFeedsWithInterval("your_stock_token", "your_interval")

	// Keep the program running
	select {}
}
