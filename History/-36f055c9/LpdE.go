package main

import "testing"

var tests = []struct {
	name     string
	dividend float32
	divison  float32
	expected float32
	isErr    bool
}

func TestDivide(t *testing.T) {
	_, err := divide(10.0, 1.0)

	if err != nil {
		t.Error("Got an error when we should have")
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := divide(10.0, 0.1)

	if err == nil {
		t.Error("Did not get an error when we should not have")
	}
}
