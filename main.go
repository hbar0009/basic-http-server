package main

import (
	"log"
	// "sockets/client"
	"sockets/server"
	// "time"
)

func main() {
	log.Print("Starting server")
	server.RunServer()

	// log.Print("Starting client")
	// go client.RunClient()

	// time.Sleep(3 * time.Second)
}
