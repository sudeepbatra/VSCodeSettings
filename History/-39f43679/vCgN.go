package main

import "log"

func main() {
	var firstLine = "Once upon a time"

	for i, l := range firstLine {
		log.Println(i, ":", l)
	}
}
