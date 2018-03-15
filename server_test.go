package main

import "testing"

func StartTest(t *testing.T) {
	StartServer("config.json")
	println("Hello world")
}
