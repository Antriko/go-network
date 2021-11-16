// TODO put everything from client.go to menu.go

package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type menuSettings struct {
	camera rl.Camera

	playerModel         userModel
	playerPos           rl.Vector3
	playerScale         rl.Vector3
	playerRotationAxis  rl.Vector3
	playerRotationAngle float32
	playerTint          rl.Color
}

func initMenu() *menuSettings {
	menu := &menuSettings{}
	menu.camera = rl.Camera{}
	menu.camera.Position = rl.NewVector3(0.0, 14.0, 16.0)
	menu.camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	menu.camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	menu.camera.Fovy = 75.0
	// rl.SetCameraMode(menu.camera, rl.CameraOrbital)

	menu.playerPos = rl.NewVector3(0.0, -16.0, -10.0)
	menu.playerScale = rl.NewVector3(5.0, 5.0, 5.0)
	menu.playerRotationAxis = rl.NewVector3(0.0, 1.0, 0.0)
	menu.playerRotationAngle = float32(0.0)
	menu.playerTint = rl.White

	menu.playerModel = userModel{}
	return menu
}

func displayMenu(buttons []button) {
	rl.UpdateCamera(&menu.camera) // camera.update(dt)
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	for i := 0; i < len(buttons); i++ {
		posY := int32(100 + (150 * i))
		rect := rl.Rectangle{X: float32(buttons[i].posX), Y: float32(posY), Width: float32(buttons[i].width), Height: float32(buttons[i].height)}
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), rect) {
			rl.DrawRectangle(buttons[i].posX, posY, buttons[i].width, buttons[i].height, rl.LightGray)
			buttons[i].isHover = true
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				buttons[i].function()
			}
		} else {
			rl.DrawRectangle(buttons[i].posX, posY, buttons[i].width, buttons[i].height, rl.DarkGray)
			buttons[i].isHover = false
		}
		rl.DrawText(buttons[i].text, buttons[i].posX+buttons[i].offsetX, posY+25, buttons[i].height/2, rl.White)
	}

	rl.BeginMode3D(menu.camera)
	menu.showPlayer()
	rl.EndMode3D()
}

func (menu *menuSettings) showPlayer() {
	// rl.DrawGrid(10, 1.0)
	menu.playerRotationAngle += 1.25
	// log.Println(menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale)
	rl.DrawModelEx(menu.playerModel.bottom, menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale, menu.playerTint)
	rl.DrawModelEx(menu.playerModel.body, menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale, menu.playerTint)
	rl.DrawModelEx(menu.playerModel.head, menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale, menu.playerTint)
	rl.DrawModelEx(menu.playerModel.hair, menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale, menu.playerTint)
	rl.DrawModelEx(menu.playerModel.accessory, menu.playerPos, menu.playerRotationAxis, menu.playerRotationAngle, menu.playerScale, menu.playerTint)
}
