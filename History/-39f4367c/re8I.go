package main

import "log"

func main() {
	myMap := make(map[string]string)

	myMap["dog"] = "Tuffy"

	log.Println("myMap['dog']:", myMap["dog"])

}
