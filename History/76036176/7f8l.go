package angelone_smartapi

import (
	"fmt"
	"time"

	SmartApi "github.com/angel-one/smartapigo"
	"github.com/angel-one/smartapigo/websocket"
)

var socketClient *websocket.SocketClient

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	err := socketClient.Subscribe()
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// Triggered when a message is received
func onMessage(message []map[string]interface{}) {
	fmt.Printf("Message Received :- %v\n", message)
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d\n", attempt)
}

func AngelOneWebSocket() {

	// Create New Angel Broking Client
	ABClient := SmartApi.New("S1632585", "4321", "iGKWS2zU")

	fmt.Println("Client :- ", ABClient)

	// User Login and Generate User Session
	session, err := ABClient.GenerateSession("960425")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	session.UserSessionTokens, err = ABClient.RenewAccessToken(session.RefreshToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("User Session Tokens :- ", session.UserSessionTokens)

	//Get User Profile
	session.UserProfile, err = ABClient.GetUserProfile()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("User Profile :- ", session.UserProfile)
	fmt.Println("User Session Object :- ", session)

	// New Websocket Client
	socketClient = websocket.New(session.ClientCode, session.FeedToken, "nse_cm|17963&nse_cm|3499&nse_cm|11536&nse_cm|21808&nse_cm|317")

	// Assign callbacks
	socketClient.OnError(onError)
	socketClient.OnClose(onClose)
	socketClient.OnMessage(onMessage)
	socketClient.OnConnect(onConnect)
	socketClient.OnReconnect(onReconnect)
	socketClient.OnNoReconnect(onNoReconnect)

	// Start Consuming Data
	socketClient.Serve()
}