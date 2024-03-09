package main

import "log"

func main() {
	animals := []string{"dog", "fish", "horse", "cow", "cat"}

	for _, animal := range animals {
		log.Println(animal)
	}
}
