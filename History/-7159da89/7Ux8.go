package main

type Person struct {
	FirstName: string `json:"first_name"`
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
