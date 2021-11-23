// TODO put everything from client.go to menu.go

package client

import (
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
	menu.rotation = 0.0

	menu.menus = make(map[string]menuButtons)
	menu.chosenModel = make(map[string]int)
	menu.chosenModel["accessory"] = 0
	menu.chosenModel["hair"] = 0
	menu.chosenModel["head"] = 0
	menu.chosenModel["body"] = 0
	menu.chosenModel["bottom"] = 0

	menu.initMenuButtons()

	// raygui.LoadGuiStyle("client/menuStyle/solarized.style") // Hover not working?
	raygui.SetStyleProperty(raygui.GlobalTextFontsize, 35)
	return menu
}

func (menu *menuSettings) displayMenu(menuName string) {
	var getMenu menuButtons
	for key, value := range menu.menus {
		if key == menuName {
			getMenu = value
		}
	}
	dt := rl.GetFrameTime()
	menu.camera.Update(dt)
	rl.BeginDrawing()

	posY := float32(100.0)
	for _, value := range getMenu.buttons {
		posX := float32(rl.GetScreenWidth()/4) - float32(value.width/2) // center on first 1/4th
		label := rl.NewRectangle(posX-float32(value.width), posY, float32(value.width), float32(value.height))
		rect := rl.NewRectangle(posX, posY, float32(value.width), float32(value.height))
		switch value.typeOf {
		case "click":
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
			raygui.Label(label, value.text)
			menu.chosenModel[strings.ToLower(value.text)] = raygui.Spinner(rect, menu.chosenModel[strings.ToLower(value.text)], 0, len(arr)-1)

			raygui.Label(rl.NewRectangle(posX+float32(value.width), posY, float32(value.width), float32(value.height)), arr[menu.chosenModel[strings.ToLower(value.text)]].name)
			if old != menu.chosenModel[strings.ToLower(value.text)] {
				value.function()
			}
		}
		posY += 120
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
	offsetX  int32
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
	btn.offsetX = 100
	btn.width = rl.MeasureText(text, btn.height/2) + (btn.offsetX * 2)
	btn.posX = int32(rl.GetScreenWidth()/2) - btn.width

	buttons = append(buttons, btn)

	return buttons
}

func (menu *menuSettings) initMenuButtons() {
	// Main menu buttons ------
	var buttons []button
	buttons = addButton(buttons, "click", "play", play)
	buttons = addButton(buttons, "click", "random", menu.selectRandomModels)
	buttons = addButton(buttons, "click", "customise", customise)
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
	// log.Println(buttons)
	menu.menus["customise"] = menuButtons{
		buttons,
		"customise",
		true,
	}
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
	for key := range models {
		var keys []string
		for mapKey := range models[key] {
			keys = append(keys, mapKey)
		}
		randomSelection := keys[rand.Intn(len(keys))]
		switch key {
		case "accessory":
			menu.playerModel.accessory = models[key][randomSelection]
		case "hair":
			menu.playerModel.hair = models[key][randomSelection]
		case "head":
			menu.playerModel.head = models[key][randomSelection]
		case "body":
			menu.playerModel.body = models[key][randomSelection]
		case "bottom":
			menu.playerModel.bottom = models[key][randomSelection]
		}
	}
}
