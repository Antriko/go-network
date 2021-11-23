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
	modelFolder := "models"
	arrayOfModels = make(map[string][]modelEntity)
	files, _ := ioutil.ReadDir(modelFolder)
	for _, f := range files {
		innerFiles, _ := ioutil.ReadDir(modelFolder + "/" + f.Name())
		models[f.Name()] = make(map[string]modelEntity)
		// -- Make first in list
		switch f.Name() {
		case "hair":
			bald := modelEntity{
				"bald",
				"hair",
				rl.Model{},
			}
			models["hair"]["bald"] = bald
			arrayOfModels["hair"] = append(arrayOfModels["hair"], bald)
		case "accessory":
			nothing := modelEntity{
				"nothing",
				"accessory",
				rl.Model{},
			}
			models["accessory"]["nothing"] = nothing
			arrayOfModels["accessory"] = append(arrayOfModels["nothing"], nothing)
		}
		//
		for _, inner := range innerFiles {
			path := modelFolder + "/" + f.Name() + "/" + inner.Name()
			ext := filepath.Ext(path)
			if ext == ".glb" {
				// Load model dynamically
				name := inner.Name()[:len(inner.Name())-len(ext)] // trim .glb from file name
				model := rl.LoadModel(path)
				modelEnt := modelEntity{
					name, // TODO Change from camelCasing to a regular word
					f.Name(),
					model,
				}
				arrayOfModels[f.Name()] = append(arrayOfModels[f.Name()], modelEnt)
				models[f.Name()][name] = modelEnt
			}
		}
	}

	// Random selection of clothing
	menu.selectRandomModels()

	for !rl.WindowShouldClose() {

		switch player.state {
		case "menu":
			menu.displayMenu("main")
		case "customise":
			menu.displayMenu("customise")

		case "game":
			game()
		}

		debugging()

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
