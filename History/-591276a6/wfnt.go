package main

import "fmt"

func main() {
	var conferenceName = "Go Conference"
	const conferenceTickets = 50
	var remainingTickets = conferenceTickets

	fmt.Printf("Welcome to", conferenceName, "booking application!")
	fmt.Println("We have a total of", conferenceTickets, "tickets and only", remainingTickets, "left")
	fmt.Println("Get your ticker here to attend")

}