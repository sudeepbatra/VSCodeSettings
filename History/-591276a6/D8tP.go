package main

import "fmt"

func main() {
	var conferenceName string = "Go Conference"
	const conferenceTickets int = 50
	var remainingTickets int = conferenceTickets

	fmt.Printf("Welcome to %v booking application!\n", conferenceName)
	fmt.Println("We have a total of", conferenceTickets, "tickets and only", remainingTickets, "left")
	fmt.Println("Get your ticker here to attend")

	var firstName string
	var lastName string
	var email string
	var userTickets int
	// ask for their name
	fmt.Println("Please enter your first name")
	fmt.Scan(&firstName)

	fmt.Println("Please enter your last name")
	fmt.Scan(&lastName)

	fmt.Println("Please enter your email")
	fmt.Scan(&email)

	fmt.Println("Please enter the number of tickets you want to buy")
	fmt.Scan(&userTickets)

	userTickets = 2
	fmt.Println(userName)
	fmt.Println(userTickets)
}
