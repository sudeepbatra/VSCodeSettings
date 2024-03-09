package main

func main() {
	type user struct {
		FirstName string
		LastName  string
		email     string
		Age       int
	}

	var users []User

	user = append(users, user("John", "Smith", "asdas", 30))
}
