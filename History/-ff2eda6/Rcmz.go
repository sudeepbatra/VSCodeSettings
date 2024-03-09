package breeze

import socketio "github.com/googollee/go-socket.io"

type SocketEventBreeze struct {
	namespace      string
	breeze         *Breeze // Assuming you have a Breeze struct defined elsewhere
	sio            *socketio.Client
	tokenList      map[string]bool
	ohlcState      map[string]string
	authentication bool
}

func NewSocketEventBreeze(namespace string, breezeInstance *Breeze) *SocketEventBreeze {
	sio, _ := socketio.NewClient([]string{"your_socket_io_server_url"}, nil)
	client := &SocketEventBreeze{
		namespace:      namespace,
		breeze:         breezeInstance,
		sio:            sio,
		tokenList:      make(map[string]bool),
		ohlcState:      make(map[string]string),
		authentication: true,
	}
	client.sio.On("connect_error", client.myConnectError)
	client.sio.On("disconnect", client.onDisconnect)
	return client
}
