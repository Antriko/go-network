package server

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func serverUserConnect() {
	log.SetFlags(log.Lshortfile)
	log.Println("User connect up")
	addr := net.UDPAddr{
		Port: 8079,
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
		log.Println(client, p)
		var tmp userInfo
		err = json.Unmarshal(p, &tmp)
		if err != nil {
			log.Println("Json error", err)
		}
		go appendUser(tmp)
		log.Println(tmp)
	}
}
func appendUser(user userInfo) {
	tmp := false
	for i := 0; i < len(userCoords); i++ {
		if userCoords[i].Username == user.Username {
			tmp = true
		}
	}
	if !tmp {
		tmpCoord := coords{}
		tmpCoord.Username = user.Username
		tmpCoord.X, tmpCoord.Y, tmpCoord.Z = 0.0, 0.0, 0.0
		tmpCoord.LastActivity = time.Now()
		userCoords = append(userCoords, tmpCoord)
	} else {
		log.Println("User already exists") // continue instance?
	}
	log.Println(userCoords)
}

type userInfo struct {
	Username string
}
