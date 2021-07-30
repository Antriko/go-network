package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type playerInfo struct {
	username      string
	pos           rl.Vector3 // 2D space for now
	movementSpeed float32
	size          rl.Vector3
	state         string
	//Texture       rl.Texture2D
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
	player.movementSpeed = 0.1

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
}

func (player *playerInfo) renderPlayer() {
	rl.DrawCubeWires(player.pos, player.size.X, player.size.Y, player.size.Z, rl.Black)
}
func (player *playerInfo) renderPlayerTag() {
	cubeScreenPosition := rl.GetWorldToScreen(rl.NewVector3(player.pos.X, player.pos.Y, player.pos.Z), camera.Camera)
	header := player.username
	rl.DrawText(header, (int32(cubeScreenPosition.X) - (rl.MeasureText(header, 100) / 2)), int32(cubeScreenPosition.Y), 20, rl.Black)
}
