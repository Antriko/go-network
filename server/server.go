package server

import (
	"encoding/json"
	"log"
	"net"
)

func Start() {
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
		p := make([]byte, 256) // make new byte every time, rather than overriding (when placed before for)
		_, client, err := server.ReadFromUDP(p)
		if err != nil {
			log.Println("Network error", err)
		}
		log.Printf("CLIENT %v : %s\n", client.Port, p)
		log.Printf("%+v", p)

		for i := 0; i < len(p); i++ {
			if p[i] == 0 {
				p = p[:i]
				break
			}
		}
		log.Println(p)

		var bytes packetInfo
		jserr := json.Unmarshal(p, &bytes)
		log.Println(json.Valid(p))
		if jserr != nil {
			log.Println("json err", jserr)
		}
		log.Printf("%+v", bytes)

		byte2string := string(p)
		switch byte2string {
		case "connect":
			log.Println("connect things")
		case "disconnect":
			log.Println("disconnect things")
		}
		go response(server, client)

	}
}

func response(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("SERVER MESSAGE"), addr)
	if err != nil {
		log.Println(err)
	}
}

type packetInfo struct {
	Username string
	Message  string
}
