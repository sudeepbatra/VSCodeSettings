package main

import (
	"fmt"
	"net/http"

	"github.com/username/learning-go-web-app/pkg/handlers"
)

const portNumber = ":8080"

// main is the main applicaton function
func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
