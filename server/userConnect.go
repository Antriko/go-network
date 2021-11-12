package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func serverUserConnect() {
	log.SetFlags(log.Lshortfile)
	addr := net.TCPAddr{
		Port: 8079,
		IP:   net.ParseIP("localhost"),
	}
	server, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Println("Network error", err)
	}
	log.Printf("User server up at port %d", addr.Port)
	for {
		p := make([]byte, 1024) // make new byte every time, rather than overriding (when placed before for)
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Println("Network error", err)
			return
		}

		go func() {
			for {
				n, err := conn.Read(p)
				usr := findUserInfo(conn, p, n)

				log.Println(string(p[:n]))

				if err != nil {
					conn.Close()
					log.Println("User disconnect", usr.Username, err)
					return
				}
				for _, value := range userConnectionsMap { // key, value
					fmt.Fprintf(value.Conn, string(p[:n]))
				}
			}
		}()

	}
}

func findUserInfo(conn *net.TCPConn, p []byte, n int) userConnInfo {
	var userInfo userInfo
	_ = json.Unmarshal(p[:n], &userInfo)

	if value, ok := userConnectionsMap[conn]; ok {
		return value
	} else {
		log.Println("Adding new user")
		// No user
		userConnected := userConnInfo{
			conn,
			userInfo.Username,
		}
		userConnectionsMap[conn] = userConnected
		return userConnected
	}
}

type userInfo struct { // TODO Maybe add user customisation ?? Clothing
	Info     string
	Username string
	Time     time.Time
}

type userConnInfo struct {
	Conn     *net.TCPConn
	Username string
	// userInfo userInfo ??
}
