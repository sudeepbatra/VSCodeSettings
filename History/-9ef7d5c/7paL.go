package main

import "errors"

func main() {

}

func divide(x, y float32) (float32, error) {
	var result float32

	if y == 0 {
		return result, errors.New("cannot divide by 0")
	}

	result = x / y
	return result, nil
}
