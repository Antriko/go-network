package main

import (
	"os"

	"github.com/Antriko/go-network/client"
	"github.com/Antriko/go-network/server"
)

func main() {
	// go run . client/server
	switch os.Args[1] {
	case "client":
		client.Start()
	case "server":
		server.Start()
	}
}
