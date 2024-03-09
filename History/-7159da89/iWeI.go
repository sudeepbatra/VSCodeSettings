package main

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HairColor string `json:"hair_color"`
	HasDog    string `json:"has_dog"`
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

}
