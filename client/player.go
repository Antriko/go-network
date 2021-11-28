package client

import (
	"net"
	"time"

	"github.com/Antriko/go-network/shared"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type playerInfo struct {
	username       string
	ID             uint32
	pos            rl.Vector3 // 2D space for now
	movementSpeed  float32
	size           rl.Vector3
	scale          float32
	state          string
	status         string
	gameStatus     string
	chatMessage    string
	serverConn     *net.TCPConn
	chatServerConn *net.TCPConn
	conn           *shared.DualConnection
	//Texture       rl.Texture2D
	model              userModel
	UserModelSelection shared.UserModelSelection
}

type userModel struct {
	accessory modelEntity
	hair      modelEntity
	head      modelEntity
	body      modelEntity
	bottom    modelEntity
}
type modelEntity struct {
	name     string
	category string
	model    rl.Model
}

func initPlayer(username string) *playerInfo {
	player := &playerInfo{}
	// 3D
	//player.pos = rl.NewVector3(0.0, 1.0, 0.0)
	//player.Texture = playerTexture

	player.state = "menu"
	player.username = username
	player.pos = rl.NewVector3(0.0, 0.0, 0.0)
	player.size = rl.NewVector3(2.0, 2.0, 2.0)
	player.scale = 2.5
	player.movementSpeed = 0.1

	player.gameStatus = "move"
	player.chatMessage = ""
	// move, typing

	// player.model = rl.LoadModel("models/castle.obj")
	player.UserModelSelection = shared.UserModelSelection{
		Accessory: 0,
		Hair:      0,
		Head:      0,
		Body:      0,
		Bottom:    0,
	}

	return player

}

func (player *playerInfo) playerTyping() {
	input := rl.GetKeyPressed()

	if input != 0 {
		// log.Println("Input: ", string(input), input)
	}
	// 259 = backspace
	if input == 259 && len(player.chatMessage) > 0 {
		player.chatMessage = player.chatMessage[:len(player.chatMessage)-1]
	}

	// A-Z, space and 0-9 for now
	if (65 <= input && input <= 90) || input == 32 || (48 <= input && input <= 57) {
		player.chatMessage = player.chatMessage + string(input)
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		if len(player.chatMessage) > 0 {
			player.sendChatMessage()
		}
		player.gameStatus = "move"
	}
}

func (player *playerInfo) playerMovement() {

	// Left - Right
	if rl.IsKeyDown(rl.KeyA) {
		player.pos.X -= player.movementSpeed
	} else if rl.IsKeyDown(rl.KeyD) {
		player.pos.X += player.movementSpeed
	}

	// Up - Down
	if rl.IsKeyDown(rl.KeyW) {
		player.pos.Z -= player.movementSpeed
	} else if rl.IsKeyDown(rl.KeyS) {
		player.pos.Z += player.movementSpeed
	}

	// Camera movement
	dt := rl.GetFrameTime()
	if rl.IsKeyDown(rl.KeyRight) { // Rotate
		camera.Angle.X += camera.RotationSpeed * dt
	} else if rl.IsKeyDown(rl.KeyLeft) {
		camera.Angle.X -= camera.RotationSpeed * dt
	}

	// Enter chat
	if rl.IsKeyPressed(rl.KeyEnter) {
		player.gameStatus = "chat"
	}
}

func (player *playerInfo) renderPlayer() {
	rl.DrawCubeWires(player.pos, player.size.X, player.size.Y, player.size.Z, rl.Black)

	rl.DrawModel(player.model.accessory.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.hair.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.head.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.body.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.bottom.model, player.pos, player.scale, rl.White)

}
func (player *playerInfo) renderPlayerTag() {
	cubeScreenPosition := rl.GetWorldToScreen(rl.NewVector3(player.pos.X, player.pos.Y, player.pos.Z), camera.Camera)
	header := player.username
	tagSize := int32(20)
	rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, tagSize) / 2)), int32(cubeScreenPosition.Y)-40, tagSize, rl.Black)
}

func (player *playerInfo) sendChatMessage() {

	DataWriteChan <- &shared.C2SChatMessagePacket{
		Username: player.username,
		Type:     shared.AllChat,
		Message:  player.chatMessage,
	}
	player.chatMessage = ""

}

type ChatMessage struct {
	Username string
	Type     shared.ChatType
	Time     time.Time
	Message  string
}
