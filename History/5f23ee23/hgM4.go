package main

import "github.com/tsawler/myniceprogram/helpers"

const numPool = 10

func CalculateValue(intChan chan int) {
	randomNumber := helpers.RandomNumber(numPool)
	intChan <- randomNumber
}

// Means of sending information from one part of the program to another part of the program.
// Unique to go
func main() {
	intChan := make(chan int)
	defer close(intChan)
}
