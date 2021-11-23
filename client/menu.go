// TODO put everything from client.go to menu.go

package client

import (
	"math"
	"math/rand"
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type menuSettings struct {
	camera *CustomCamera

	playerModel userModel
	playerPos   rl.Vector3
	playerScale float32
	playerTint  rl.Color

	rotation float32

	cameraMove cameraMove

	menus       map[string]menuButtons
	chosenModel map[string]int
}

type menuButtons struct {
	buttons    []button
	menuTitle  string
	showPlayer bool // Show player in menu
}

func initMenu() *menuSettings {
	menu := &menuSettings{}
	menu.camera = NewCustomCamera(10.0, 128.0, 8.0)
	menu.camera.SetTarget(rl.NewVector3(0.0, 3.0, 0.0))
	menu.camera.Angle = rl.NewVector2(-0.0, -0.4)

	menu.playerPos = rl.NewVector3(8.0, -4.0, 0.0)
	scale := float32(6.0)
	menu.playerScale = scale
	menu.playerTint = rl.White

	menu.playerModel = userModel{}
	menu.rotation = 4.0

	menu.menus = make(map[string]menuButtons)
	menu.chosenModel = make(map[string]int)
	menu.chosenModel["accessory"] = 0
	menu.chosenModel["hair"] = 0
	menu.chosenModel["head"] = 0
	menu.chosenModel["body"] = 0
	menu.chosenModel["bottom"] = 0

	menu.initMenuButtons()

	return menu
}

func (menu *menuSettings) displayMenu(menuName string) {
	posY := float32(rl.GetScreenHeight()) * 0.05                        // 1/20th of the top
	width := float32(rl.GetScreenWidth()) * 0.2                         // 1/5th of screen width
	height := float32(math.Round(float64(rl.GetScreenWidth()) * 0.025)) // yeah.
	labelTextSize := height / 2
	raygui.SetStyleProperty(raygui.GlobalTextFontsize, int64(height)) // dynamically set font size depending on screen height
	var getMenu menuButtons
	for key, value := range menu.menus {
		if key == menuName {
			getMenu = value
		}
	}
	dt := rl.GetFrameTime()
	menu.camera.Update(dt)
	rl.BeginDrawing()

	for _, value := range getMenu.buttons {
		posX := float32(rl.GetScreenWidth()/4) - float32(width/2) // center on first 1/4th

		rect := rl.NewRectangle(posX, posY, width, height)
		switch value.typeOf {
		case "button":
			buttonClicked := raygui.Button(rect, value.text)
			if buttonClicked {
				value.function()
			}
		case "spinner":
			var arr []modelEntity
			for key, models := range arrayOfModels {
				if key == strings.ToLower(value.text) {
					arr = models
				}
			}
			old := menu.chosenModel[strings.ToLower(value.text)]
			menu.chosenModel[strings.ToLower(value.text)] = raygui.Spinner(rect, menu.chosenModel[strings.ToLower(value.text)], 0, len(arr)-1)
			if old != menu.chosenModel[strings.ToLower(value.text)] {
				value.function()
			}
			// Labels
			// raygui.Label won't allow for font size change.
			rl.DrawText(value.text, int32(posX-float32(rl.MeasureText(value.text, int32(labelTextSize))))-int32(labelTextSize), int32(posY), int32(labelTextSize), rl.Black)
			rl.DrawText(arr[menu.chosenModel[strings.ToLower(value.text)]].name, int32(posX+width)+int32(labelTextSize), int32(posY), int32(labelTextSize), rl.Black)
		}
		posY += height * 1.5
	}

	rl.ClearBackground(rl.RayWhite)

	if getMenu.showPlayer {
		offTop := int32(10)
		rl.DrawLine(int32(rl.GetScreenWidth()/2), offTop, int32(rl.GetScreenWidth()/2), int32(rl.GetScreenHeight())-offTop, rl.Black)
		rl.BeginMode3D(menu.camera.Camera)
		menu.showPlayer()
		menu.controls()
		rl.EndMode3D()
	}
}

type button struct {
	typeOf   string
	text     string
	function func()
	isHover  bool
	posX     int32
	width    int32
	height   int32
}

func addButton(buttons []button, typeOf string, text string, function func()) []button {
	// TODO dynamic button width or text center?
	btn := button{}
	btn.typeOf = typeOf
	btn.text = text
	btn.function = function
	btn.isHover = false
	btn.height = 100
	btn.width = rl.MeasureText(text, btn.height/2) + (100 * 2)
	btn.posX = int32(rl.GetScreenWidth()/2) - btn.width

	buttons = append(buttons, btn)

	return buttons
}

func (menu *menuSettings) initMenuButtons() {
	// Main menu buttons ------
	var buttons []button
	buttons = addButton(buttons, "button", "play", play)
	buttons = addButton(buttons, "button", "random", menu.selectRandomModels)
	buttons = addButton(buttons, "button", "customise", customise)
	menu.menus["main"] = menuButtons{
		buttons,
		"main",
		true,
	}
	buttons = nil

	// Customisation menu buttions ------
	// log.Println("ok")
	// for key, value := range models {
	// 	log.Println(key)
	// 	for mapKey := range value {
	// 		log.Println(key, mapKey)
	// 	}
	// }

	buttons = addButton(buttons, "spinner", "Accessory", menu.changeModel)
	buttons = addButton(buttons, "spinner", "Hair", menu.changeModel)
	buttons = addButton(buttons, "spinner", "Head", menu.changeModel)
	buttons = addButton(buttons, "spinner", "Body", menu.changeModel)
	buttons = addButton(buttons, "spinner", "Bottom", menu.changeModel)
	buttons = addButton(buttons, "button", "Random", menu.selectRandomModels)
	buttons = addButton(buttons, "button", "Save", backToMainMenu)
	// log.Println(buttons)
	menu.menus["customise"] = menuButtons{
		buttons,
		"customise",
		true,
	}
}

func backToMainMenu() {
	player.state = "menu"
}

func (menu *menuSettings) showPlayer() {
	menu.playerModel.body.model.Transform = rl.MatrixRotateY(menu.rotation)
	menu.playerModel.hair.model.Transform = rl.MatrixRotateY(menu.rotation)
	menu.playerModel.head.model.Transform = rl.MatrixRotateY(menu.rotation)
	menu.playerModel.accessory.model.Transform = rl.MatrixRotateY(menu.rotation)
	menu.playerModel.bottom.model.Transform = rl.MatrixRotateY(menu.rotation)
	rl.DrawModel(menu.playerModel.bottom.model, menu.playerPos, menu.playerScale, menu.playerTint)
	rl.DrawModel(menu.playerModel.body.model, menu.playerPos, menu.playerScale, menu.playerTint)
	rl.DrawModel(menu.playerModel.head.model, menu.playerPos, menu.playerScale, menu.playerTint)
	rl.DrawModel(menu.playerModel.hair.model, menu.playerPos, menu.playerScale, menu.playerTint)
	rl.DrawModel(menu.playerModel.accessory.model, menu.playerPos, menu.playerScale, menu.playerTint)
}

type cameraMove struct {
	original      float32
	originalMouse int32 // X axis only
	newMouse      int32
	started       bool
}
type rotationStruct struct {
	axis  rl.Vector3
	angle float32
}

func (menu *menuSettings) controls() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && int(rl.GetMouseX()) > rl.GetScreenWidth()/2 {
		menu.cameraMove = cameraMove{
			menu.rotation,
			rl.GetMouseX(),
			rl.GetMouseX(),
			true,
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && menu.cameraMove.started {
		menu.cameraMove.newMouse = rl.GetMouseX()

		differenceBetweenMouse := float32(menu.cameraMove.originalMouse - menu.cameraMove.newMouse)
		differenceBetweenMouse = differenceBetweenMouse / menu.camera.RotationSpeed

		menu.rotation = menu.cameraMove.original + float32(differenceBetweenMouse)
	}

	// Disable ability to rotate player when mouse is on left side
	if rl.IsMouseButtonUp(rl.MouseLeftButton) && int(rl.GetMouseX()) < rl.GetScreenWidth()/2 {
		menu.cameraMove.started = false
	}

}

func (menu *menuSettings) selectRandomModels() {
	for key, value := range arrayOfModels {
		var keys []string
		for mapKey := range models[key] {
			keys = append(keys, mapKey)
		}
		randNum := rand.Intn(len(keys))
		randomSelection := value[randNum]
		switch key {
		case "accessory":
			menu.playerModel.accessory = randomSelection
			menu.chosenModel["accessory"] = randNum
		case "hair":
			menu.playerModel.hair = randomSelection
			menu.chosenModel["hair"] = randNum
		case "head":
			menu.playerModel.head = randomSelection
			menu.chosenModel["head"] = randNum
		case "body":
			menu.playerModel.body = randomSelection
			menu.chosenModel["body"] = randNum
		case "bottom":
			menu.playerModel.bottom = randomSelection
			menu.chosenModel["bottom"] = randNum
		}
	}
}
