package main

import (
	"encoding/json"
	"log"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HairColor string `json:"hair_color"`
	HasDog    bool   `json:"has_dog"`
}

func main() {
	myJson := `
[
    {
            "first_name": "Wally",
            "last_name": "West",
            "hair_color": "red",
            "has_dog": false
    },
    {
            "first_name": "Wally",
            "last_name": "West",
            "hair_color": "red",
            "has_dog": false
    }
]`

	var unmarshalled []Person

	err := json.Unmarshal([]byte(myJson), &unmarshalled)

	if err != nil {
		log.Println("Error unmarshalling the json", err)
	}

	log.Printf("unmarshalled: %v", unmarshalled)

	//write json from a struct

	var mySlice []Person
	var m1 Person
	m1.FirstName = "asdfa"
	m1.LastName = "sadkfj"
	m1.HairColor = "sdlkfjdls"
	m1.HasDog = true
}
