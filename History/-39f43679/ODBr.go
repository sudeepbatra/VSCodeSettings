package main

import "log"

func main() {
	animals := []string{"dog", "fish", "horse", "cow", "cat"}

	for i, animal := range animals {
		log.Println(i, animal)
	}
}
