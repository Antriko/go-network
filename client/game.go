package client

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/Antriko/go-network/shared"
	"github.com/Antriko/go-network/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var players = make(map[string]shared.Coords)
var connectedPlayers = make(map[string]shared.OtherPlayer)
var chatHistory []ChatMessage
var chatScroll = 0

func game() {
	log.SetFlags(log.Lshortfile)

	switch player.gameStatus {
	case "move":
		player.playerMovement()
	case "chat":
		player.playerTyping()
	case "genMesh":
		worldMap.Instances = make(map[world.TileType]world.Instance)
		worldMap.MeshGen(worldOption)
		player.gameStatus = "move"
	}

	camera.SetTarget(player.pos)
	camera.SetPosition(player.pos.X, player.pos.Y+10, player.pos.Z+10)
	player.playerRotation()

	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(camera.Camera)
	renderGridFloor()

	worldMap.RenderMesh()

	player.renderPlayer()
	renderOthers()
	rl.EndMode3D()
	player.renderPlayerTag()
	renderOtherTag()
	renderChat()

}

type coords struct {
	Username     string
	X            float32
	Y            float32
	Z            float32
	LastActivity time.Time
}

func renderGridFloor() {
	// grid of light and dark green
	slices := 10
	spacing := float32(5.0)
	// -splice*spacing to +splice*spacing in XZ
	// use modulus%2 to interchange between colours, add +x for every z
	for x := slices * -1; x < slices; x++ {
		for z := slices * -1; z < slices; z++ {
			// Would be a lot cleaner if ternary operators worked in go
			// col = (z+x)%2 == 0 ? dark : light
			var col color.RGBA
			if (z+x)%2 == 0 {
				col = rl.NewColor(2, 35, 28, 255)
			} else {
				col = rl.NewColor(0, 77, 37, 255)
			}
			rl.DrawCube(rl.NewVector3(float32(x)*spacing, 0, float32(z)*spacing), spacing, 0, spacing, col)
		}
	}

}

func renderOthers() {
	for key, value := range players {
		if key != player.username {
			playerPos := rl.NewVector3(value.X, value.Y, value.Z)
			// rl.DrawCubeWires(playerPos, player.Size.X, player.Size.Y, player.Size.Z, rl.Black) // default size - players
			renderFromUserModelSelection(connectedPlayers[value.Username].UserModelSelection, playerPos, value.Facing)
		}
	}
}

func renderFromUserModelSelection(UserModelSelection shared.UserModelSelection, pos rl.Vector3, rotation float32) {
	arrayOfModels["accessory"][UserModelSelection.Accessory].model.Transform = rl.MatrixRotateY(rotation * (math.Pi / 180))
	arrayOfModels["hair"][UserModelSelection.Accessory].model.Transform = rl.MatrixRotateY(rotation * (math.Pi / 180))
	arrayOfModels["head"][UserModelSelection.Accessory].model.Transform = rl.MatrixRotateY(rotation * (math.Pi / 180))
	arrayOfModels["body"][UserModelSelection.Accessory].model.Transform = rl.MatrixRotateY(rotation * (math.Pi / 180))
	arrayOfModels["bottom"][UserModelSelection.Accessory].model.Transform = rl.MatrixRotateY(rotation * (math.Pi / 180))

	rl.DrawModel(arrayOfModels["accessory"][UserModelSelection.Accessory].model, pos, player.scale, rl.White)
	rl.DrawModel(arrayOfModels["hair"][UserModelSelection.Hair].model, pos, player.scale, rl.White)
	rl.DrawModel(arrayOfModels["head"][UserModelSelection.Head].model, pos, player.scale, rl.White)
	rl.DrawModel(arrayOfModels["body"][UserModelSelection.Body].model, pos, player.scale, rl.White)
	rl.DrawModel(arrayOfModels["bottom"][UserModelSelection.Bottom].model, pos, player.scale, rl.White)
}

func renderOtherTag() {
	for key, value := range players {
		if key != player.username {
			tagOffset := 6
			cubeV3 := rl.NewVector3(value.X, value.Y+float32(tagOffset), value.Z)
			cubeScreenPosition := rl.GetWorldToScreen(cubeV3, camera.Camera)
			header := value.Username
			difference := rl.Vector3Subtract(rl.NewVector3(player.pos.X, player.pos.Y, player.pos.Z), cubeV3).Z
			if difference < 30 { // Don't load if too far
				tmpDiff := difference
				if difference > 20 {
					tmpDiff = 20
				}
				tagSize := int32(math.Abs(20 - float64(tmpDiff))) // neg to pos, pos to neg
				rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, tagSize) / 2)), int32(cubeScreenPosition.Y), tagSize, rl.Black)
			}

		}
	}
}

func renderChat() {
	var chatLen int
	fontSize := 20
	if len(chatHistory) > 6 {
		chatLen = len(chatHistory) - 6
	} else {
		chatLen = 0
	}

	largestText := int32(rl.GetScreenWidth() / 3)
	for i := len(chatHistory) - 1; i >= chatLen; i-- {
		message := getMessage(chatHistory[i])
		current := rl.MeasureText(message, int32(fontSize)) + int32(fontSize)
		if current > largestText {
			largestText = current
		}
	}

	rl.DrawRectangle(0, int32(rl.GetScreenHeight()-(fontSize*2)-(fontSize*len(chatHistory))), largestText, int32(fontSize*len(chatHistory))+(int32(fontSize)*2), rl.ColorAlpha(rl.DarkGray, .5))

	for i := len(chatHistory) - 1; i >= chatLen; i-- {
		// log.Println(chatHistory[i].Time.Format("15:04:05"))
		message := getMessage(chatHistory[i])
		rl.DrawText(message, 10, int32(rl.GetScreenHeight()-fontSize-(fontSize*(len(chatHistory)-i))), int32(fontSize), rl.Black)
	}

	rl.DrawLine(10, int32(rl.GetScreenWidth()), 150, int32(rl.GetScreenWidth())+5, rl.Black)
	message := fmt.Sprintf("%s: %s", player.username, player.chatMessage)
	rl.DrawText(message, 10, int32(rl.GetScreenHeight()-fontSize), int32(fontSize), rl.Black)
	if player.gameStatus == "chat" { // blinking |
		if time.Now().Nanosecond() > 500000000 {
			rl.DrawText("|", int32(fontSize)+rl.MeasureText(message, int32(fontSize)), int32(rl.GetScreenHeight()-fontSize), int32(fontSize), rl.Black)
		}
	}
}

func getMessage(message ChatMessage) string {
	switch message.Type {
	case shared.AllChat: // All chat
		return fmt.Sprintf("[%v] %s: %s", message.Time.Format("15:04:05"), message.Username, message.Message)
	case shared.UserConnect: // User connected
		return fmt.Sprintf("[%v] User %s has connected.", message.Time.Format("15:04:05"), message.Username)
	case shared.UserDisconnect:
		return fmt.Sprintf("[%v] User %s has disonnected.", message.Time.Format("15:04:05"), message.Username)
	case shared.CommandWorldSize: // User changed map size
		return fmt.Sprintf("[%v] User %s has changed world size to %s.", message.Time.Format("15:04:05"), message.Username, message.Message)
	}
	log.Println(message.Type)
	return "?"
}

// Gets Y pos of tile that player is in
// TODO - 3x3 collision instead of returning y value of noise?
func getYCollision(pos rl.Vector3) float32 {

	offset := float32(0)
	if worldMap.Size%2 == 1 {
		offset = (worldOption.TileSize / 2) + float32((worldMap.Size-int(worldOption.TileSpacing))/2)
	} else {
		offset = float32(worldMap.Size) / 2
	}

	posX := float32(pos.X/worldOption.TileSpacing + offset)
	posZ := float32(pos.Z/worldOption.TileSpacing + offset)

	if posX < 0 {
		posX = 0
	} else if posX > float32(worldMap.Size) {
		posX = float32(worldMap.Size) - 1
	}
	if posZ < 0 {
		posZ = 0
	} else if posZ > float32(worldMap.Size) {
		posZ = float32(worldMap.Size) - 1
	}

	worldTile := worldMap.Tiles[int(math.Floor(float64(posZ)))][int(math.Floor(float64(posX)))]

	if worldTile.Noise < 0.3 {
		return 0.3 * worldOption.TileSize * worldOption.HeightMulti
	} else {
		return float32(worldTile.Noise*float64(worldOption.TileSize*worldOption.HeightMulti)) + (worldOption.TileSize / 2)
	}
}
