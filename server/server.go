package server

import (
	"github.com/Antriko/go-network/shared"
)

// Store all connections within one global
var AllConnections Connections

// Coords of all users
var userCoordsMap = make(map[string]shared.Coords) // UDP
type Connections map[*shared.DualConnection]*UserConnection

type UserConnection struct {
	Connection         *shared.DualConnection
	Username           string
	ID                 uint32
	UserModelSelection shared.UserModelSelection
	// Add more information if needed
}

func Start() { // goroutine for all the for{} loop
	go server()
	for {
	}
}
