package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

func serverChat() {
	log.SetFlags(log.Lshortfile)
	addr := net.TCPAddr{
		Port: 8070,
		IP:   net.ParseIP("localhost"),
	}
	server, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Println("Network error", err)
	}
	log.Printf("Chat server up at port %d", addr.Port)
	// Handle connections

	for {
		p := make([]byte, 1024)
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Println("TCP Conn err", err)
			return
		}
		go func() {
			for {
				n, err := conn.Read(p)
				chatConnInit(conn, p, n)
				if err != nil {
					return
				}

				// send message to all other connections
				for key := range chatConnectionsMap { // key, value
					fmt.Fprintf(key, string(p[:n]))
				}
			}
		}()

	}

}
func chatConnInit(conn *net.TCPConn, p []byte, n int) {
	if _, ok := userConnectionsMap[conn]; !ok {
		chatConnectionsMap[conn] = userConnInfo{
			conn,
			string(p[:n]),
		}
	}
}

type ChatMessage struct {
	Info     string
	username string
	message  string
	time     time.Time
}
