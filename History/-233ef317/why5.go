package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	gosocketio "github.com/ambelovsky/gosf-socketio"
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
	socket       *gosocketio.Client
	socketOrder  *gosocketio.Client
	socketOHLC   *gosocketio.Client
	userId       string
	sessionToken string
	apiKey       string
	config       map[string]string // Assuming you have a config map
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func (sm *SocketManager) connect(isOrder, isOHLC bool) {
	authDict := map[string]string{
		"user":   sm.userId,
		"token":  sm.sessionToken,
		"appkey": sm.apiKey,
	}

	if isOrder && sm.socketOrder == nil {
		options := gosocketio.Options{
			Transport: "websocket",
			Query:     authDict,
		}
		uriOrder := sm.config["LIVE_FEEDS_URL"]
		sm.socketOrder = gosocketio.NewClient(strings.TrimSuffix(uriOrder, "/")+"/socket.io/", options)
		sm.socketOrder.On("connect", func(h *gosocketio.Channel, args Message) {
			fmt.Println("Connected to Order Socket.IO")
		})
		sm.socketOrder.Connect()
	}

	if isOHLC && sm.socketOHLC == nil {
		options := gosocketio.Options{
			Transport: "websocket",
			Query:     authDict,
		}
		uriOHLC := sm.config["LIVE_OHLC_STREAM_URL"]
		sm.socketOHLC = gosocketio.NewClient(strings.TrimSuffix(uriOHLC, "/")+"/socket.io/", options)
		sm.socketOHLC.On("connect", func(h *gosocketio.Channel, args Message) {
			fmt.Println("Connected to OHLC Socket.IO")
		})
		sm.socketOHLC.Connect()
	}

	if sm.socket == nil {
		options := gosocketio.Options{
			Transport: "websocket",
			Query:     authDict,
			Header:    http.Header{"User-Agent": []string{"node-socketio[client]/socket"}},
		}
		uri := sm.config["LIVE_STREAM_URL"]
		sm.socket = gosocketio.NewClient(strings.TrimSuffix(uri, "/")+"/socket.io/", options)
		sm.socket.On("connect", func(h *gosocketio.Channel, args Message) {
			fmt.Println("Connected to Socket.IO")
		})
		sm.socket.Connect()
	}
}

func (sm *SocketManager) connectTicker() {
	sm.connect(false, false)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sm := &SocketManager{
		userId:       "your_user_id",
		sessionToken: "your_session_token",
		apiKey:       "your_api_key",
		config: map[string]string{
			"LIVE_FEEDS_URL":       "https://livefeeds.icicidirect.com",
			"LIVE_STREAM_URL":      "https://livestream.icicidirect.com",
			"LIVE_OHLC_STREAM_URL": "https://breezeapi.icicidirect.com",
		},
	}

	sm.connectTicker()

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				go sendJoin(sm.socket)
				go sendJoin(sm.socket)
				go sendJoin(sm.socket)
				go sendJoin(sm.socket)
				go sendJoin(sm.socket)
			}

			time.Sleep(10 * time.Second)
			sm.socket.Close()
		}()
	}

	time.Sleep(6000 * time.Second)
	log.Println(" [x] Complete")
}
