package client

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Antriko/go-network/shared"
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
	}

	camera.SetTarget(player.pos)
	camera.SetPosition(player.pos.X, player.pos.Y+10, player.pos.Z+10)
	player.playerRotation()

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

}

type coords struct {
	Username     string
	X            float32
	Y            float32
	Z            float32
	LastActivity time.Time
}

func renderOthers() {
	for key, value := range players {
		if key != player.username {
			playerPos := rl.NewVector3(value.X, value.Y, value.Z)
			rl.DrawCubeWires(playerPos, player.size.X, player.size.Y, player.size.Z, rl.Black) // default size - players
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
			cubeV3 := rl.NewVector3(value.X, value.Y, value.Z)
			cubeScreenPosition := rl.GetWorldToScreen(cubeV3, camera.Camera)
			header := value.Username
			difference := rl.Vector3Subtract(rl.NewVector3(player.pos.X, player.pos.Y, player.pos.Z), cubeV3).Z
			if difference < 30 { // Don't load if too far
				tmpDiff := difference
				if difference > 20 {
					tmpDiff = 20
				}
				tagSize := int32(math.Abs(20 - float64(tmpDiff))) // neg to pos, pos to neg
				rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, tagSize) / 2)), int32(cubeScreenPosition.Y)-40, tagSize, rl.Black)
			}

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
		// log.Println(chatHistory[i].Time.Format("15:04:05"))
		var message string
		switch chatHistory[i].Type {
		case shared.AllChat: // All chat
			message = fmt.Sprintf("[%v] %s: %s", chatHistory[i].Time.Format("15:04:05"), chatHistory[i].Username, chatHistory[i].Message)
		case shared.UserConnect: // User connected
			message = fmt.Sprintf("[%v] User %s has connected.", chatHistory[i].Time.Format("15:04:05"), chatHistory[i].Username)
		case shared.UserDisconnect:
			message = fmt.Sprintf("[%v] User %s has disonnected.", chatHistory[i].Time.Format("15:04:05"), chatHistory[i].Username)
		}
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
