package server

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func serverCoords() {
	log.SetFlags(log.Lshortfile)
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("localhost"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println("Network error", err)
	}
	log.Printf("Coords server up at port %d", addr.Port)

	go func() {
		for {
			p := make([]byte, 1024)
			n, client, err := server.ReadFromUDP(p)
			if err != nil {
				log.Println("Network error", err)
			}
			var coords coords
			err = json.Unmarshal(p[:n], &coords)
			if err != nil {
				log.Println("Json error", err, client)
			}
			go coordUpdate(coords, server, client)
		}
	}()
}
func coordUpdate(userCoord coords, conn *net.UDPConn, addr *net.UDPAddr) {

	userCoordsMap[userCoord.Username] = coords{
		userCoord.Username,
		userCoord.X,
		userCoord.Y,
		userCoord.Z,
		time.Now(),
	}

	// convert from map to array
	var coordsArray []coords
	for _, value := range userCoordsMap {
		coordsArray = append(coordsArray, value)
	}
	jsonData, err := json.Marshal(coordsArray)
	if err != nil {
		log.Println(err)
	}
	_, jsonErr := conn.WriteToUDP(jsonData, addr)
	if jsonErr != nil {
		log.Println(err)
	}
	log.Println(len(userCoordsMap))
}

type coords struct {
	Username     string
	X            float32
	Y            float32
	Z            float32
	LastActivity time.Time
}
