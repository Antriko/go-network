package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/Antriko/go-network/shared"
	"github.com/Antriko/go-network/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// global
var menu *menuSettings = initMenu()
var player *playerInfo = initPlayer("tmp")
var camera = NewCustomCamera(10.0, 3.0, 100.0)
var models = make(map[string]map[string]modelEntity)
var arrayOfModels = make(map[string][]modelEntity)
var users = make(map[string]shared.OtherPlayer)

var worldMap *world.WorldStruct
var worldOption world.WorldOptionStruct
var worldSize int

// TODO move camera to game/player struct instead of global

func Start() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	rl.SetTraceLog(rl.LogError)
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "Chatroom")
	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	worldOption = world.WorldOptionStruct{
		TileSize:    5.0,
		TileSpacing: 5.0,
		HeightMulti: 2.0,
	}

	worldSize = 50
	worldMap = world.CreateWorld(worldSize)
	worldMap.MeshGen(worldOption)

	player.username = ""

	// Model loading - Load all models (dynamically, reading each file dir) and save into a global model map with nested maps
	// map["head"]map["default"] ect
	var jsonModel map[string][]modelJson
	b, _ := ioutil.ReadFile("models/models.json")
	_ = json.Unmarshal(b, &jsonModel)
	arrayOfModels = make(map[string][]modelEntity)
	for key, value := range jsonModel {
		models[key] = make(map[string]modelEntity)
		for _, value2 := range value {
			var model rl.Model
			if value2.File == "" { // For no hair and accessory
				model = rl.Model{}
			} else {
				model = rl.LoadModel(value2.File)
			}
			modelEnt := modelEntity{
				value2.Name,
				key,
				model,
			}
			arrayOfModels[key] = append(arrayOfModels[key], modelEnt)
			models[key][value2.Name] = modelEnt
		}
	}

	// Default to 0 0 0 0 0
	menu.changeModel()

	for !rl.WindowShouldClose() {

		// if rl.IsKeyPressed(rl.KeySpace) {
		// 	worldMap = world.CreateWorld(worldSize)
		// 	worldMap.MeshGen(worldOption)
		// } else if rl.IsKeyPressed(rl.KeyMinus) {
		// 	worldSize--
		// 	if worldSize < 2 {
		// 		worldSize = 2
		// 	}
		// 	worldMap = world.CreateWorld(worldSize)
		// 	worldMap.MeshGen(worldOption)
		// } else if rl.IsKeyPressed(rl.KeyEqual) {
		// 	worldSize++
		// 	worldMap = world.CreateWorld(worldSize)
		// 	worldMap.MeshGen(worldOption)
		// }

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
	// some validation before conneting to game/server
	if player.username == "" {
		return
	}
	go serverConn()
	player.state = "game"
	// create player

}

type userInfo struct {
	Info               string
	Username           string
	Time               time.Time
	UserModelSelection shared.UserModelSelection
}

type modelJson struct {
	Name string
	File string
}

func debugging() {
	// Debugging
	start := 10
	posX := int32(10)
	posY := int32(10)
	incrY := func() int {
		start += 10
		return start
	}
	rl.DrawText(fmt.Sprintf("Username: %v", player.username), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("MousePos: %v", rl.GetMousePosition()), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("playerState: %v", player.state), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("gameState: %v", player.gameStatus), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("playerPos: %v", player.pos), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("rotation: %v", menu.rotation), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("playerFACING: %v", player.rotation.facing), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("playerROTATION: %v", player.rotation.rotation), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("worldSize: %v", worldMap.Size), posX, int32(incrY()), posY, rl.Black)

}
