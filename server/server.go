package server

import "net"

// TODO
// Change from arrays to map

// global

// global maps of connections
var userCoordsMap = make(map[string]coords)                  // UDP
var userConnectionsMap = make(map[*net.TCPConn]userConnInfo) // TCP
var chatConnectionsMap = make(map[*net.TCPConn]userConnInfo) // TCP

func Start() { // goroutine for all the for{} loop
	go serverUserConnect()
	go serverChat()
	go serverCoords()
	for {
	}
}
