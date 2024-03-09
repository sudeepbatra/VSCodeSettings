package main

import "fmt"

var myName string

func main() {
	fmt.Println("Hello World!")

	var whatToSay string = "adsfsdf"
	var i int

	fmt.Println(whatToSay)

	i = 64
	fmt.Println("i is set to", i)

	whatWasSaid := saySomething()
	fmt.Println("what was said", whatWasSaid)
}

func saySomething() string {
	return "something"
}
