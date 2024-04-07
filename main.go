package main

import "golang-chat-server/network" // Import the package that contains the "network" identifier

func main() {
	server := network.NewServer()
	server.StartServer()
}
