package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// global
var menu *menuSettings = initMenu()
var player *playerInfo = initPlayer("tmp")
var camera = NewCustomCamera(10.0, 3.0, 100.0)
var models = make(map[string]map[string]modelEntity)
var arrayOfModels = make(map[string][]modelEntity)

// TODO move camera to game/player struct instead of global

func Start() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	rl.SetTraceLog(rl.LogError)
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "Chatroom")
	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	username := fmt.Sprint(rand.Intn(10000))
	player.username = username
	// Model loading - Load all models (dynamically, reading each file dir) and save into a global model map with nested maps
	// map["head"]map["default"] ect

	var jsonModel map[string]map[string]modelJson
	b, _ := ioutil.ReadFile("models/models.json")
	_ = json.Unmarshal(b, &jsonModel)
	arrayOfModels = make(map[string][]modelEntity)
	for key, value := range jsonModel {
		log.Println(key, value)
		models[key] = make(map[string]modelEntity)
		for key2, value2 := range value {
			log.Println(value2)
			var model rl.Model
			if value2.File == "" { // For no hair and accessory
				model = rl.Model{}
			} else {
				model = rl.LoadModel(value2.File)
			}
			modelEnt := modelEntity{
				value2.Name,
				key2,
				model,
			}
			arrayOfModels[key] = append(arrayOfModels[key], modelEnt)
			log.Println(key, key2)
			models[key][key2] = modelEnt
		}
	}

	// Default to 0 0 0 0 0
	menu.changeModel()

	for !rl.WindowShouldClose() {

		switch player.state {
		case "menu":
			menu.displayMenu("main")
		case "customise":
			menu.displayMenu("customise")

		case "game":
			game()
		}

		// debugging()

		rl.EndDrawing()
	}

	rl.CloseWindow()
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
}

type userInfo struct {
	Info     string
	Username string
	Time     time.Time
}

type modelJson struct {
	Name string
	File string
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
	rl.DrawText(fmt.Sprintf("rotation: %v", menu.rotation), 20, int32(incrY()), 20, rl.Black)
}
