package main

// func main() {
// 	fmt.Println("Hello World!")

// 	var whatToSay string = "adsfsdf"
// 	var i int

// 	fmt.Println(whatToSay)

// 	i = 64
// 	fmt.Println("i is set to", i)

// 	whatWasSaid, whatElseWasSaid := saySomething()
// 	fmt.Println("what was said", whatWasSaid, whatElseWasSaid)
// }

// func saySomething() (string, string) {
// 	return "something", "else"
// }

// import "log"

// func main() {
// 	var myString string
// 	myString = "Green"

// 	log.Println("myString is set to", myString)
// 	changeUsingPointer(&myString)
// 	log.Println("myString is set to", myString)
// }

// func changeUsingPointer(s *string) {
// 	log.Println("s address is set to", s)
// 	log.Println("s value is set to", *s)
// 	newValue := "Red"
// 	*s = newValue
// }

import (
	"log"
	"time"
)

// var s string
var s = "seven"

var firstName string
var lastName string
var phoneNumber string
var age int
var birthDate time.Time

type User struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Age         int
	BirthDate   time.Time
}

func main() {
	var s = "six"
	log.Println(s)

}

func saySomething(s string) (string, string) {
	return s, "World"
}
