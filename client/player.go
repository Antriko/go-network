package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type playerInfo struct {
	username      string
	pos           rl.Vector2 // 2D space for now
	movementSpeed float32
	Size          rl.Vector2
	//Texture       rl.Texture2D
}

func initPlayer(username string) *playerInfo {
	player := &playerInfo{}
	// 3D
	//player.pos = rl.NewVector3(0.0, 1.0, 0.0)
	//player.Size = rl.NewVector3(2.0, 2.0, 2.0)
	//player.Texture = playerTexture

	player.username = username
	player.pos = rl.NewVector2(0.0, 0.0)
	player.movementSpeed = 0.25

	return player
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
		player.pos.Y -= player.movementSpeed
	} else if rl.IsKeyDown(rl.KeyS) {
		player.pos.Y += player.movementSpeed
	}
}

func (player *playerInfo) renderPlayer() {
	rl.DrawRectangle(int32(player.pos.X), int32(player.pos.Y), int32(player.Size.X), int32(player.Size.Y), rl.Black)
}
