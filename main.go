package main

import (
	"golang-chat-server/network" // Import the package that contains the "network" identifier
	_ "net/http/pprof"
)

func main() {
	server := network.NewServer()
	server.StartServer()
}
