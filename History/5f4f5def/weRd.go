package main

import (
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	n, err := fmt.Fprintf(w, "Hello, world!")

	// 	if err != nil {
	// 		log.Println("Error while trying to write to the response writer", err)
	// 	}

	// 	log.Println("Bytes written: ", n)
	// })

	http.HandleFunc("/", Home)

	_ = http.ListenAndServe(":8080", nil)
}
