package main

import "fmt"

func main() {
	fmt.Println("Hello World!")

	var whatToSay string = "adsfsdf"
	var i int

	fmt.Println(whatToSay)

	i = 64
	fmt.Println("i is set to", i)

	saySomething()
}

func saySomething() string {
	return "something"
}