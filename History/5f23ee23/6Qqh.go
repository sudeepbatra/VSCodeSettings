package main

import "github.com/tsawler/myniceprogram/helpers"

func CalculateValue(intChan chan int) {
	randomNumer := helpers.RandomNumber(10)
}

// Means of sending information from one part of the program to another part of the program.
// Unique to go
func main() {
	intChan := make(chan int)

}
