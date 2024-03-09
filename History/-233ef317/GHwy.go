package main

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

type MySocketClient struct {
	userId       string
	sessionToken string
	apiKey       string
	socket       *socketio.Client
	socketOrder  *socketio.Client
	socketOHLC   *socketio.Client
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
		opts := &socketio.Options{
			Transport: []string{"websocket"},
			Query:     authDict,
		}
		uriOrder := config.urls[Config.UrlEnum.LIVE_FEEDS_URL]
		socketOrder, err := socketio.NewClient(uriOrder, opts)
		if err != nil {
			// Handle error
		}
		c.socketOrder = socketOrder
		c.socketOrder.OnConnect(func() {
			fmt.Println("Connected to order socket")
		})
		c.socketOrder.Connect()
	}

	if isOHLC && c.socketOHLC == nil {
		opts := &socketio.Options{
			Transport: []string{"websocket"},
			Query:     authDict,
			Path:      "/ohlcvstream/",
		}
		uriohlc := config.urls[Config.UrlEnum.LIVE_OHLC_STREAM_URL]
		socketOHLC, err := socketio.NewClient(uriohlc, opts)
		if err != nil {
			// Handle error
		}
		c.socketOHLC = socketOHLC
		c.socketOHLC.OnConnect(func() {
			fmt.Println("Connected to OHLC socket")
		})
		c.socketOHLC.Connect()
	} else if c.socket == nil {
		opts := &socketio.Options{
			Transport: []string{"websocket"},
			Query:     authDict,
		}
		uri := config.urls[Config.UrlEnum.LIVE_STREAM_URL]
		socket, err := socketio.NewClient(uri, opts)
		if err != nil {
			// Handle error
		}
		c.socket = socket
		c.socket.OnConnect(func() {
			fmt.Println("Connected to default socket")
		})
		c.socket.Connect()
	}
}

func (c *MySocketClient) connectTicker() {
	c.connect(false, false)
}

func (c *MySocketClient) watch(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("joining stock subscription:", stock)
			c.socket.Emit("join", stock)
		}
	}
}

func (c *MySocketClient) unWatch(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			fmt.Println("leaving stock subscription:", stock)
			c.socket.Emit("leave", stock)
		}
	}
}

func (c *MySocketClient) watchOHLC(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			c.socketOHLC.Emit("join", stock)
		}
	}
}

func (c *MySocketClient) unWatchOHLC(stocks []string) {
	if stocks != nil && len(stocks) > 0 {
		for _, stock := range stocks {
			c.socketOHLC.Emit("leave", stock)
		}
	}
}

func (c *MySocketClient) subscribeFeeds(stockToken string) (map[string]interface{}, error) {
	if stockToken == "" {
		return c.socketConnectionResponse(config.responseMessage[Config.ResponseEnum.BLANK_STOCK_CODE]), nil
	}

	c.watch([]string{stockToken})
	return c.socketConnectionResponse(fmt.Sprintf(config.responseMessage[Config.ResponseEnum.STOCK_SUBSCRIBE_MESSAGE], stockToken)), nil
}

func (c *MySocketClient) subscribeFeedsWithInterval(stockToken, interval string) (map[string]interface{}, error) {
	if interval == "" || !contains(config.typeLists[Config.ListEnum.INTERVAL_TYPES_STREAM_OHLC], interval) {
		return nil, fmt.Errorf(config.exceptionMessage[Config.ExceptionEnum.STREAM_OHLC_INTERVAL_ERROR])
	}

	c.interval = config.channelIntervalMap[interval]
	c.connect(false, true)

	if stockToken == "" {
		return c.socketConnectionResponse(config.responseMessage[Config.ResponseEnum.BLANK_STOCK_CODE]), nil
	}

	c.watchOHLC([]string{stockToken})
	return c.socketConnectionResponse(fmt.Sprintf(config.responseMessage[Config.ResponseEnum.STOCK_SUBSCRIBE_MESSAGE], stockToken)), nil
}

func (c *MySocketClient) socketConnectionResponse(message string) map[string]interface{} {
	// Implement your response logic here
	// You can return a JSON response as needed
	response := make(map[string]interface{})
	response["message"] = message
	return response
}

func (c *MySocketClient) errorException(errMsg string) {
	// Implement your error handling logic here
	fmt.Println("Error:", errMsg)
}

func contains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func main() {
	// Replace with your actual credentials
	userId := "your_user_id"
	sessionToken := "your_session_token"
	apiKey := "your_api_key"

	client := NewMySocketClient(userId, sessionToken, apiKey)
	// Use the client to connect and perform other operations
}
