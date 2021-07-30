package server

// TODO
// Add last time active and remove user if inactive for certain amount of time

// global
var userCoords []coords

func Start() { // goroutine for all except last
	go serverUserConnect()
	serverCoords()
}
