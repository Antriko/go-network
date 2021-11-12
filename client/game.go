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
var chatHistory []ChatMessage
var chatScroll = 0

func game() {
	log.SetFlags(log.Lshortfile)

	switch player.gameStatus {
	case "move":
		player.playerMovement()
	case "chat":
		player.playerTyping()
	}

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
	renderChat()
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
	log.Println(string(jsonData))

	fmt.Fprintf(conn, "%s", jsonData) // string vulnerability or something
	n, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
	p = p[0:n]

	var playerCoords []coords
	err = json.Unmarshal(p, &playerCoords)
	if err != nil {
		//log.Printf("JSON ERR: %v", err)
		return
	}
	players = playerCoords
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

func renderChat() {

	var chatLen int
	if len(chatHistory) > 6 {
		chatLen = len(chatHistory) - 6
	} else {
		chatLen = 0
	}
	for i := len(chatHistory) - 1; i >= chatLen; i-- {
		message := fmt.Sprintf("[%v:%v:%v] %s: %s", chatHistory[i].Time.Hour(), chatHistory[i].Time.Minute(), chatHistory[i].Time.Second(), chatHistory[i].Username, chatHistory[i].Message)
		rl.DrawText(message, 10, int32(rl.GetScreenHeight()-20-(20*(len(chatHistory)-i))), 20, rl.Black)
	}

	rl.DrawLine(10, int32(rl.GetScreenWidth()), 150, int32(rl.GetScreenWidth())+5, rl.Black)
	message := fmt.Sprintf("%s: %s", player.username, player.chatMessage)
	rl.DrawText(message, 10, int32(rl.GetScreenHeight())-20, 20, rl.Black)
	if player.gameStatus == "chat" { // blinking |
		if time.Now().Nanosecond() > 500000000 {
			rl.DrawText("|", 20+rl.MeasureText(message, 20), int32(rl.GetScreenHeight())-20, 20, rl.Black)
		}
	}
}
