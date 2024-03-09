package main

func main() {
	type User struct {
		FirstName string
		LastName  string
		email     string
		Age       int
	}

	var users []User

	users = append(users, User{"John", "Smith", "asdas", 30})
	users = append(users, User{"asdasd", "asdasd", "asdas", 30343})
}
