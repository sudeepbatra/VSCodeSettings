package breeze

type SocketEventBreeze struct {
	namespace      string
	breeze         *Breeze // Assuming you have a Breeze struct defined elsewhere
	sio            *socketio.Client
	tokenList      map[string]bool
	ohlcState      map[string]string
	authentication bool
}
