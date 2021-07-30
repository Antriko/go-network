// TODO put everything from client.go to menu.go

package client

import rl "github.com/gen2brain/raylib-go/raylib"

func menu(buttons []button) {
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
		rl.DrawText(buttons[i].text, buttons[i].posX+10, posY+25, buttons[i].height/2, rl.White)
	}
}
