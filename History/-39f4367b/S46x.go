package main

import "log"

func main() {
	myNum := 100
	isTrue := false

	if myNum && isTrue {
		log.Println("myNum is greater than 99 and isTrue is set to true")
	}
}