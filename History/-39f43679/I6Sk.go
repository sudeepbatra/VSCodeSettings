package main

import "log"

func main() {
	animals := make(map[string]string)
	animals["dog"] = "Fido"
	animals["cat"] = "Fluffy"
	// animals := []string{"dog", "fish", "horse", "cow", "cat"}

	for _, animal := range animals {
		log.Println(animal)
	}
}
