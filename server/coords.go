package server

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func serverCoords() {
	log.SetFlags(log.Lshortfile)
	log.Println("Coords Server up", userCoords)
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("localhost"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println("Network error", err)
	}
	for {
		p := make([]byte, 1024) // make new byte every time, rather than overriding (when placed before for)
		n, client, err := server.ReadFromUDP(p)
		if err != nil {
			log.Println("Network error", err)
		}
		p = p[0:n]
		var tmp coords
		err = json.Unmarshal(p, &tmp)
		if err != nil {
			log.Println("Json error", err, client)
		}
		go coordUpdate(tmp, server, client)
	}
}
func coordUpdate(userCoord coords, conn *net.UDPConn, addr *net.UDPAddr) {
	for i := 0; i < len(userCoords); i++ {
		if userCoords[i].Username == userCoord.Username {
			userCoords[i].X = userCoord.X
			userCoords[i].Y = userCoord.Y
			userCoords[i].Z = userCoord.Z
			userCoords[i].LastActivity = time.Now()
		}
	}
	jsonEncoded, err := json.Marshal(userCoords)
	if err != nil {
		log.Println(err)
	}
	_, jsonErr := conn.WriteToUDP([]byte(jsonEncoded), addr)
	if jsonErr != nil {
		log.Println(err)
	}
}

type coords struct {
	Username     string
	X            float32
	Y            float32
	Z            float32
	LastActivity time.Time
}
