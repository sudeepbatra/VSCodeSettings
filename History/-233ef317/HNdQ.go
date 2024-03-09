package main

type MySocketClient struct {
	userId       string
	sessionToken string
	apiKey       string
	socket       *socketio.Client
}
