package server

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

// TODO
// Add last time active and remove user if inactive for certain amount of time

// global
var users []coordsPacket // where all user coords are stored
var userCoords []coords

func Start() {
	go serverUserConnect()
	go serverCoords()

	log.SetFlags(log.Lshortfile)
	log.Println("Server start")
	addr := net.UDPAddr{
		Port: 8080,
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

		// remove buffer
		p = p[0:n]

		var decodedJson coordsPacket
		jserr := json.Unmarshal(p, &decodedJson)
		if jserr != nil {
			log.Println("json err", jserr)
		}

		log.Printf("CLIENT %v : %s\n", decodedJson.Username, p)
		switch decodedJson.Message {
		case "connect":
			log.Println("connect things")
			go connectUser(decodedJson, server, client)
		case "test":
			go testResponse(server, client)
		case "coords":
			go updateCoords(decodedJson, server, client)
		}
	}
}

func updateCoords(data coordsPacket, conn *net.UDPConn, addr *net.UDPAddr) {
	for i := 0; i < len(users); i++ {
		if users[i].Username == data.Username {
			users[i].CoordX = data.CoordX
			users[i].CoordY = data.CoordY
			users[i].CoordZ = data.CoordZ
		}
	}

	jsonEncoded, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
	}
	_, jsonErr := conn.WriteToUDP([]byte(jsonEncoded), addr)
	if jsonErr != nil {
		log.Println(err)
	}
}
func connectUser(data coordsPacket, conn *net.UDPConn, addr *net.UDPAddr) {
	tmp := false
	for i := 0; i < len(users); i++ {
		if users[i].Username == data.Username {
			tmp = true
		}
	}
	if !tmp {
		users = append(users, data)
	}
}
func testResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	tme := time.Now()
	_, err := conn.WriteToUDP([]byte(tme.String()), addr)
	if err != nil {
		log.Println(err)
	}
}

type coordsPacket struct {
	Username string
	Message  string
	CoordX   float32
	CoordY   float32
	CoordZ   float32
}
