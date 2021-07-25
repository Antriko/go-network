package main

import (
	"os"

	"github.com/Antriko/go-network/client/client.go"
	"github.com/Antriko/go-network/server/server.go"
)

func main() {
	switch os.Args[1] {
	case "client":
		client.Start()
	case "server":
		server.Start()
	}
}
