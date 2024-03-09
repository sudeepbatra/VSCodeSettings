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

	mySlice = append(mySlice, m1)

	var m2 Person
	m2.FirstName = "dfsdsdfa"
	m2.LastName = "sdf"
	m2.HairColor = "hfgh"
	m2.HasDog = true

	mySlice = append(mySlice, m2)

	newJson, err := json.MarshalIndent(mySlice, "", "	")

	if err != nil {
		log.Println("Error in Marshalling the struct", err)
	}

	log.Println(string(newJson))
}
