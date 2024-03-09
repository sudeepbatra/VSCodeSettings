package main

type myStruct struct {
	FirstName string
}

func main() {
	var myVar myStruct
	myVar.FirstName = "John"

	myVar2 := myStruct{
		FirstName: "Mary",
	}
}
