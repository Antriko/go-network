package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// global
var menu *menuSettings = initMenu()
var player *playerInfo = initPlayer("tmp")
var camera = NewCustomCamera(10.0, 3.0, 100.0)
var models = make(map[string]map[string]rl.Model)

func Start() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	username := fmt.Sprint(rand.Intn(10000))
	player.username = username

	screenWidth := int32(1280)
	screenHeight := int32(720)
	rl.InitWindow(screenWidth, screenHeight, "Multiplayer")

	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	// Model loading - Load all models (dynamically, reading each file dir) and save into a global model map with nested maps
	// map["head"]map["default"] ect
	modelFolder := "models"
	files, _ := ioutil.ReadDir(modelFolder)
	for _, f := range files {
		innerFiles, _ := ioutil.ReadDir(modelFolder + "/" + f.Name())
		models[f.Name()] = make(map[string]rl.Model)
		for _, inner := range innerFiles {
			path := modelFolder + "/" + f.Name() + "/" + inner.Name()
			ext := filepath.Ext(path)
			if ext == ".glb" {
				log.Println(path)
				// Load model dynamically
				name := inner.Name()[:len(inner.Name())-len(ext)]
				model := rl.LoadModel(path)
				models[f.Name()][name] = model // trim .glb from file name
			}
		}
	}

	for key := range models {
		var keys []string
		for mapKey := range models[key] {
			keys = append(keys, mapKey)
		}
		log.Println(keys)
		randomSelection := keys[rand.Intn(len(keys))]
		switch key {
		case "accessory":
			menu.playerModel.accessory = models[key][randomSelection]
		case "hair":
			menu.playerModel.hair = models[key][randomSelection]
		case "head":
			menu.playerModel.head = models[key][randomSelection]
		case "body":
			menu.playerModel.body = models[key][randomSelection]
		case "bottom":
			menu.playerModel.bottom = models[key][randomSelection]
		}
	}

	var buttons []button
	buttons = addButton(buttons, "play", play)

	rl.SetCameraMode(camera.Camera, rl.CameraCustom) // Set a first person misc.CustomCamera mode
	dt := rl.GetFrameTime()                          // delta time
	camera.SetTarget(rl.NewVector3(0.0, 0.0, 0.0))
	camera.Update(dt)

	rl.SetCameraMode(menu.camera, rl.CameraFree)
	for !rl.WindowShouldClose() {

		switch player.state {
		case "menu":
			displayMenu(buttons)

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
	go player.connectUser()
	go player.connChat()
	player.state = "game"
	// create player

}
func (player *playerInfo) connectUser() {

	addr := net.TCPAddr{
		Port: 8079,
		IP:   net.ParseIP("localhost"),
	}
	conn, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		log.Println("Chat server error", err)
	}
	player.serverConn = conn
	jsonData, err := json.Marshal(userInfo{
		"connect",
		player.username,
		time.Now(),
	})
	if err != nil {
		log.Printf("JSON data err %v", err)
		return
	}
	fmt.Fprintf(player.serverConn, "%s", jsonData) // send user init to server

	// handle server response
	p := make([]byte, 1024)
	go func() {
		for {
			n, err := player.serverConn.Read(p)
			if err != nil {
				log.Println("TCP data err", err)
			}
			log.Printf("%s: %s", conn.RemoteAddr().String(), string(p[:n]))
			// add/remove users

			var userInfo userInfo
			err = json.Unmarshal(p[:n], &userInfo)
			log.Println(userInfo)
			// User connected/disconnected info
			chatHistory = append(chatHistory, ChatMessage{
				"message",
				"[SERVER]",
				"User " + userInfo.Username + " has " + userInfo.Info + ".",
				userInfo.Time,
			})

		}
	}()

	// p := make([]byte, 1024)
	// conn, err := net.Dial("udp", "localhost:8079")
	// if err != nil {
	// 	log.Printf("Some error %v", err)
	// 	return
	// }
	// jsonData, err := json.Marshal(userInfo{player.username})

	// if err != nil {
	// 	log.Printf("Some error %v", err)
	// 	return
	// }

	// fmt.Fprintf(conn, "%s", jsonData) // strings vulnerability or something
	// _, err = bufio.NewReader(conn).Read(p)
	// // TODO
	// // Wait for server response (user exists, user continue?) and then go into game
	// if err == nil {
	// 	fmt.Printf("%s\n", p)
	// } else {
	// 	fmt.Printf("Some error %v\n", err)
	// }
	// //conn.Close()
}

type userInfo struct {
	Info     string
	Username string
	Time     time.Time
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
	rl.DrawText(fmt.Sprintf("gameState: %v", player.gameStatus), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("playerPos: %v", player.pos), 20, int32(incrY()), 20, rl.Black)
	rl.DrawText(fmt.Sprintf("chatMessage: %v", player.chatMessage), 20, int32(incrY()), 20, rl.Black)
}
