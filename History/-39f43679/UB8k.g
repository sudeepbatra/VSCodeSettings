package main

import "log"

func main() {
	type User struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
	}

	var users []User

	users = append(users, User{"John", "Smith", "asdas", 30})
	users = append(users, User{"asdasd", "asdasd", "asdas", 30343})

	for _, l := range users {
		log.Println(l)
	}
}
