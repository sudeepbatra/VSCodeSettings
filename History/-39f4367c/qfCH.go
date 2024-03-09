package main

import "log"

func main() {
	var mySlice []string

	mySlice = append(mySlice, "Sudeep")
	mySlice = append(mySlice, "Batra")

	log.Println(mySlice)
}
