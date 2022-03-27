package client

import (
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
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
	rotation           playerRotation
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
type playerRotation struct {
	rotation      float32
	rotationSpeed float32
	facing        float32
	timeCount     float32
}

func initPlayer(username string) *playerInfo {
	player := &playerInfo{}

	// random spawn spot
	// between -xy to xy
	max := 10
	min := -10

	player.state = "menu"
	player.username = username
	rand.Seed(time.Now().UnixNano()) // new seed so players dont have same pos
	player.pos = rl.NewVector3(float32(rand.Intn(max-min)+min), 0.0, float32(rand.Intn(max-min)+min))
	player.size = rl.NewVector3(2.0, 2.0, 2.0)
	player.scale = 2.5
	player.movementSpeed = 0.1

	// random facing when spawn
	// set to 180 if want to face down when spawn
	player.rotation.facing = float32(rand.Intn(7) * 45)
	player.rotation.rotation = player.rotation.facing
	player.rotation.rotationSpeed = 0.25

	player.gameStatus = "move"
	player.chatMessage = ""
	// move, typing

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
	// Also / for command
	if (65 <= input && input <= 90) || input == 32 || (48 <= input && input <= 57) || input == 47 {
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

	// Determine where player is facing
	left, right, up, down := false, false, false, false

	// Left - Right
	if rl.IsKeyDown(rl.KeyA) {
		player.pos.X -= player.movementSpeed
		left = true
	} else if rl.IsKeyDown(rl.KeyD) {
		player.pos.X += player.movementSpeed
		right = true
	}

	// Up - Down
	if rl.IsKeyDown(rl.KeyW) {
		player.pos.Z -= player.movementSpeed
		up = true
	} else if rl.IsKeyDown(rl.KeyS) {
		player.pos.Z += player.movementSpeed
		down = true
	}

	// Get where player is facing
	switch {
	case up && left:
		player.rotation.facing = 315
		player.rotation.timeCount = 0
		break
	case up && right:
		player.rotation.facing = 45
		player.rotation.timeCount = 0
		break
	case down && left:
		player.rotation.facing = 225
		player.rotation.timeCount = 0
		break
	case down && right:
		player.rotation.facing = 135
		player.rotation.timeCount = 0
		break
	case up:
		player.rotation.facing = 0
		player.rotation.timeCount = 0
		break
	case right:
		player.rotation.facing = 90
		player.rotation.timeCount = 0
		break
	case down:
		player.rotation.facing = 180
		player.rotation.timeCount = 0
		break
	case left:
		player.rotation.facing = 270
		player.rotation.timeCount = 0
		break
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

	player.pos.Y = getYCollision(player.pos)

}

func (player *playerInfo) playerRotation() {
	// pivot the rotation towards facing
	// slow down once eaching facing

	// To make sure distance is between 0 and 360
	clockwiseDistance := normaliseAngle(player.rotation.facing - player.rotation.rotation)
	antiClockwiseDistance := normaliseAngle(player.rotation.rotation - player.rotation.facing)

	// Securely set threshold to complete lerp once lerp transition is close enough
	threshold := float32(5.0)
	if clockwiseDistance < threshold || antiClockwiseDistance < threshold {
		player.rotation.rotation = player.rotation.facing
		return
	}

	rotation := rl.NewVector2(player.rotation.rotation, 0.0)
	facing := rl.NewVector2(player.rotation.facing, 0.0)

	// Find quickest then set facing
	if clockwiseDistance < antiClockwiseDistance {
		if player.rotation.facing < player.rotation.rotation {
			facing.X += 360
		}
	} else {
		if player.rotation.facing > player.rotation.rotation {
			facing.X -= 360
		}
	}

	newRotation := rl.Vector2Lerp(rotation, facing, player.rotation.rotationSpeed)
	player.rotation.rotation = normaliseAngle(newRotation.X)
}
func normaliseAngle(angle float32) float32 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

func (player *playerInfo) renderPlayer() {

	player.model.accessory.model.Transform = rl.MatrixRotateY(player.rotation.rotation * (math.Pi / 180))
	player.model.hair.model.Transform = rl.MatrixRotateY(player.rotation.rotation * (math.Pi / 180))
	player.model.head.model.Transform = rl.MatrixRotateY(player.rotation.rotation * (math.Pi / 180))
	player.model.body.model.Transform = rl.MatrixRotateY(player.rotation.rotation * (math.Pi / 180))
	player.model.bottom.model.Transform = rl.MatrixRotateY(player.rotation.rotation * (math.Pi / 180))

	// rl.DrawCubeWires(player.pos, player.Size.X, player.Size.Y, player.Size.Z, rl.Black)

	rl.DrawModel(player.model.accessory.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.hair.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.head.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.body.model, player.pos, player.scale, rl.White)
	rl.DrawModel(player.model.bottom.model, player.pos, player.scale, rl.White)

	// boxHeight := player.size.Y * player.size.Z
	// rl.DrawCubeWires(rl.NewVector3(player.pos.X, player.pos.Y+(boxHeight/2), player.pos.Z), player.size.X, boxHeight, player.size.Z, rl.Red)
}
func (player *playerInfo) renderPlayerTag() {
	tagOffset := 6
	cubeV3 := rl.NewVector3(player.pos.X, player.pos.Y+float32(tagOffset), player.pos.Z)
	cubeScreenPosition := rl.GetWorldToScreen(cubeV3, camera.Camera)
	header := player.username

	tagSize := int32(20)
	rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, tagSize) / 2)), int32(cubeScreenPosition.Y)-int32(tagOffset), tagSize, rl.Black)
}

func (player *playerInfo) sendChatMessage() {

	// Check if command
	if player.chatMessage[0:1] == "/" {
		chat := strings.Split(player.chatMessage, " ")

		command := strings.ToUpper(chat[0])
		switch command {
		case "/SIZE":

			// Make sure there is a 2nd argument and that it's a integer and more than 2
			i, err := strconv.Atoi(chat[1])
			if err != nil || len(chat) < 2 || i < 2 {
				log.Println(err)
				break
			}

			DataWriteChan <- &shared.C2SChatMessagePacket{
				Username: player.username,
				Type:     shared.CommandWorldSize,
				Message:  strconv.Itoa(i),
			}
		}
	} else {
		DataWriteChan <- &shared.C2SChatMessagePacket{
			Username: player.username,
			Type:     shared.AllChat,
			Message:  player.chatMessage,
		}
	}

	player.chatMessage = ""

}

type ChatMessage struct {
	Username string
	Type     shared.ChatType
	Time     time.Time
	Message  string
}
