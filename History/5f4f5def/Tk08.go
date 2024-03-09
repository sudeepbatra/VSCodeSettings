package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {

}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {

}

// addValues adds two integers and returns the sum
// func addValues(x, y int) int {
// 	return x + y
// }

// func Divide(w http.ResponseWriter, r *http.Request) {
// 	f, err := divideValues(100.0, 0.0)

// 	if err != nil {
// 		fmt.Fprintf(w, "cannot divide by 0")
// 		return
// 	}

// 	fmt.Fprintf(w, fmt.Sprintf("%f divided by %f is %f", 100.0, 0.0, f))
// }

// func divideValues(x, y float32) (float32, error) {
// 	if y <= 0 {
// 		err := errors.New("cannot divide by zero")
// 		return 0, err
// 	}
// 	result := x / y
// 	return result, nil
// }

// main is the main applicaton function
func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	n, err := fmt.Fprintf(w, "Hello, world!")

	// 	if err != nil {
	// 		log.Println("Error while trying to write to the response writer", err)
	// 	}

	// 	log.Println("Bytes written: ", n)
	// })

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/divide", Divide)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
