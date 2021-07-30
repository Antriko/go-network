package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var players []coords

func game() {
	log.SetFlags(log.Lshortfile)
	player.playerMovement()
	camera.SetTarget(player.pos)
	camera.SetPosition(player.pos.X, player.pos.Y+10, player.pos.Z+10)

	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(camera.Camera)
	rl.DrawGrid(10, 1.0)
	player.renderPlayer()
	renderOthers()
	rl.EndMode3D()
	player.renderPlayerTag()
	renderOtherTag()
	go sendServerStatus()

}

func sendServerStatus() {
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", "localhost:8081")
	if err != nil {
		log.Printf("UDP ERR: %v", err)
		return
	}

	jsonData, err := json.Marshal(coords{player.username, player.pos.X, player.pos.Y, player.pos.Z, time.Now()})

	if err != nil {
		log.Printf("JSON ERR: %v", err)
		return
	}

	fmt.Fprintf(conn, "%s", jsonData) // string vulnerability or something
	n, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
	p = p[0:n]

	var tmp []coords
	err = json.Unmarshal(p, &tmp)

	if err != nil {
		log.Printf("JSON ERR: %v", err)
		return
	}
	players = tmp
}

type coords struct {
	Username     string
	X            float32
	Y            float32
	Z            float32
	LastActivity time.Time
}

func renderOthers() {
	for i := 0; i < len(players); i++ {
		if players[i].Username != player.username { // dont render self
			playerPos := rl.NewVector3(players[i].X, players[i].Y, players[i].Z)
			rl.DrawCubeWires(playerPos, player.size.X, player.size.Y, player.size.Z, rl.Black) // default size - players
		}
	}
}
func renderOtherTag() {
	for i := 0; i < len(players); i++ {
		if players[i].Username != player.username { // dont render self
			cubeScreenPosition := rl.GetWorldToScreen(rl.NewVector3(players[i].X, players[i].Y, players[i].Z), camera.Camera)
			header := players[i].Username
			rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, 100) / 2)), int32(cubeScreenPosition.Y), 20, rl.Black)
		}
	}
}
