package main

import "fmt"

func main() {
	var conferenceName string = "Go Conference"
	const conferenceTickets uint = -50
	var remainingTickets int = conferenceTickets

	fmt.Printf("Welcome to %v booking application!\n", conferenceName)
	fmt.Println("We have a total of", conferenceTickets, "tickets and only", remainingTickets, "left")
	fmt.Println("Get your ticker here to attend")

	var userName string
	var userTickets int
	// ask for their name

	userName = "John Doe"
	userTickets = 2
	fmt.Println(userName)
	fmt.Println(userTickets)
}
