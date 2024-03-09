package main

// func main() {
// 	fmt.Println("Hello World!")

// 	var whatToSay string = "adsfsdf"
// 	var i int

// 	fmt.Println(whatToSay)

// 	i = 64
// 	fmt.Println("i is set to", i)

// 	whatWasSaid, whatElseWasSaid := saySomething()
// 	fmt.Println("what was said", whatWasSaid, whatElseWasSaid)
// }

// func saySomething() (string, string) {
// 	return "something", "else"
// }

import "log"

func main() {
	var myString string
	myString = "Green"

	log.Println("myString is set to", myString)
}

func changeUsingPointer(s *string) {

}
