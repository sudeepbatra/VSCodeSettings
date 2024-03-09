package main

import (
	"fmt"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Hello, world!")

		if err != nil {
			log.Println("Error while trying to write to the response writer", err)
		}

		log.Println("Bytes written: ", n)
	})

	_ = http.ListenAndServe(":8080", nil)
}
