package main

const local string = "LOCAL"

func main() {
	_ = StartServer("config.json")
	select {}
}
