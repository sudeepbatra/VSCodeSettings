package main

import (
	socketio "github.com/googollee/go-socket.io"
)

type MySocketClient struct {
	userId       string
	sessionToken string
	apiKey       string
	socket       *socketio.Client
}