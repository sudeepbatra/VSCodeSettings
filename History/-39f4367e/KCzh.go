package main

import "log"

type myStruct struct {
	FirstName string
}

func (m *myStruct) printFirstName() string {

}

func main() {
	var myVar myStruct
	myVar.FirstName = "John"

	myVar2 := myStruct{
		FirstName: "Mary",
	}

	log.Println("myVar is set to", myVar.FirstName)
	log.Println("myVar2 is set to", myVar2.FirstName)
}
