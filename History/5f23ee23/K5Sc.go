package main

import "log"

// Means of sending information from one part of the program to another part of the program.
// Unique to go
func main() {
	PrintText("Hi")
}

func PrintText(s string) {
	log.Println(s)
}
