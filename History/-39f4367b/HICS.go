package main

import "log"

func main() {
	cat := "cat2"

	if cat == "cat" {
		log.Println("Cat is cat")
	} else {
		log.Println("Cat is not cat")
	}
}
