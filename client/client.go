package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// global
var player *playerInfo = initPlayer("tmp")
var camera = NewCustomCamera(10.0, 3.0, 100.0)

func Start() {
	rand.Seed(time.Now().UnixNano())
	username := fmt.Sprint(rand.Intn(10000))
	player.username = username

	screenWidth := int32(1280)
	screenHeight := int32(720)
	rl.InitWindow(screenWidth, screenHeight, "Multiplayer")
	rl.SetTargetFPS(60)

	var buttons []button
	buttons = addButton(buttons, "play", play)

	rl.SetCameraMode(camera.Camera, rl.CameraCustom) // Set a first person misc.CustomCamera mode
	dt := rl.GetFrameTime()                          // delta time
	camera.SetTarget(rl.NewVector3(0.0, 0.0, 0.0))
	camera.Update(dt)

	for !rl.WindowShouldClose() {

		switch player.state {
		case "menu":
			menu(buttons)
		case "game":
			game()
		}

		debugging()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

type button struct {
	text     string
	function func()
	isHover  bool
	posX     int32
	offsetX  int32
	width    int32
	height   int32
}

func addButton(buttons []button, text string, function func()) []button {
	// TODO dynamic button width or text center?
	btn := button{}
	btn.text = text
	btn.function = function
	btn.isHover = false
	btn.height = 100
	btn.offsetX = 100
	btn.width = rl.MeasureText(text, btn.height/2) + (btn.offsetX * 2)
	btn.posX = int32(rl.GetScreenWidth()/2) - btn.width

	buttons = append(buttons, btn)

	return buttons
}

func play() {
	go connectUser()
	player.state = "game"
	// create player

}
func connectUser() {
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", "localhost:8079")
	if err != nil {
		log.Printf("Some error %v", err)
		return
	}
	jsonData, err := json.Marshal(userInfo{player.username})

	if err != nil {
		log.Printf("Some error %v", err)
		return
	}

	fmt.Fprintf(conn, "%s", jsonData) // string vulnerability or something
	_, err = bufio.NewReader(conn).Read(p)
	// TODO
	// Wait for server response (user exists, user continue?) and then go into game
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

type userInfo struct {
	Username string
}

func debugging() {
	// Debugging
	start := 20
	incrY := func() int {
		start += 20
		return start
	}
	rl.DrawText(fmt.Sprintf("Username: %v", player.username), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("MousePos: %v", rl.GetMousePosition()), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("playerState: %v", player.state), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("playerPos: %v", player.pos), 20, int32(incrY()), 20, rl.Black)
}
