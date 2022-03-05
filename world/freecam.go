package world

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var world *worldStruct = createWorld(25)
var cam rl.Camera
var mapMatrix rl.Matrix

type worldOptionStruct struct {
	tileSize    float32
	tileSpacing float32
	heightMulti float32
}

var worldOption worldOptionStruct

func Freecam() {
	rl.InitWindow(1280, 720, "Freecam")
	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	worldOption = worldOptionStruct{5, 5, 2}
	world.meshGen()

	pos := rl.NewVector3(10, 20, 10)
	tar := rl.NewVector3(0, 0, 0)
	up := rl.NewVector3(0, 1.0, 0)
	fovy := float32(45)
	cam = rl.NewCamera3D(pos, tar, up, fovy, rl.CameraPerspective)
	rl.SetCameraMode(cam, rl.CameraFree)

	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&cam)
		if rl.IsKeyDown(rl.KeyZ) {
			cam.Target.Y += 5
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(cam)
		rl.DrawGrid(10, worldOption.tileSpacing)
		world.renderMesh()

		rl.EndMode3D()
		debugging()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

// Not efficiently rendered - deprecated
// Ignores shader rending
func (world *worldStruct) renderWorld() {
	spacing := float32(5.0)
	width := float32(world.size)
	length := width
	for y := range world.tiles {
		for x := range world.tiles[y] {

			xPos := (float32(x) - (width / 2)) * worldOption.tileSpacing
			yPos := (float32(y) - (length / 2)) * worldOption.tileSpacing

			var col color.RGBA
			// if (x+y)%2 == 0 {
			// 	col = rl.NewColor(2, 35, 28, 255)
			// } else {
			// 	col = rl.NewColor(0, 77, 37, 255)
			// }

			if world.tiles[y][x].noise < 0.15 {
				col = rl.NewColor(0, 45, 57, 255)
			} else if world.tiles[y][x].noise < 0.3 {
				col = rl.NewColor(0, 163, 204, 255)
			} else if world.tiles[y][x].noise < 0.6 {
				col = rl.NewColor(95, 99, 68, 255)
			} else {
				col = rl.NewColor(65, 69, 47, 255)
			}

			height := float32(world.tiles[y][x].noise) * worldOption.tileSize
			if height < 0.3 {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height*worldOption.tileSize, yPos)
			rl.DrawCube(pos, spacing, worldOption.tileSize, spacing, col)
			rl.DrawCubeWires(pos, spacing, worldOption.tileSize, spacing, rl.Black)
		}
	}
}

var tileInstances int
var tileTranslations []rl.Matrix // Locations of instances
var tileModel rl.Model

var treeInstances int
var treeTranslations []rl.Matrix
var treeModel rl.Model

func (world *worldStruct) meshGen() {

	// Collection of translations in which the GPU will render the same model over those positions
	// to save compute time from initialising the model over and over
	// Shader returns the colour of mesh depending on height level

	tileInstances = world.size * world.size
	tileTranslations = make([]rl.Matrix, tileInstances) // Locations of instances;

	// Create basic cube mesh
	tileModel = rl.LoadModelFromMesh(rl.GenMeshCube(worldOption.tileSpacing, worldOption.tileSize, worldOption.tileSpacing))
	tex := rl.GenImageChecked(2, 2, 1, 1, rl.Red, rl.Green)
	texture := rl.LoadTextureFromImage(tex)
	rl.SetMaterialTexture(tileModel.Materials, rl.MapDiffuse, texture)
	rl.UnloadImage(tex)

	// Shader
	tileShader := rl.LoadShader("world/glsl330/tileShader.vs", "world/glsl330/tileShader.fs")
	tileShader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(tileShader, "mvp"))
	tileShader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(tileShader, "instanceTransform"))
	tileModel.Materials.Shader = tileShader

	// Assign location for all cubes to create map
	width := float32(world.size)
	length := width
	for y := range world.tiles {
		for x := range world.tiles[y] {
			xPos := (float32(x) - (width / 2)) * worldOption.tileSpacing
			yPos := (float32(y) - (length / 2)) * worldOption.tileSpacing
			height := float32(world.tiles[y][x].noise) * worldOption.tileSize
			if height < 0.3*worldOption.tileSize {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height*worldOption.heightMulti, yPos)

			index := y*int(width) + x
			mat := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
			tileTranslations[index] = mat
			tileTranslations[index] = rl.MatrixMultiply(tileTranslations[index], rl.MatrixTranslate(0, 0, 0))

			// log.Println(index, instances, rl.Vector3Transform(rl.NewVector3(0, 0, 0), translations[index]).Y)

			// Get total amount of trees to get amount of instances needed
			if world.tiles[y][x].tile == tree {
				treeInstances++
			}
		}
	}

	// Tree model and shader init
	treeModel = rl.LoadModel("models/world/Tree.glb")
	treeShader := rl.LoadShader("world/glsl330/treeShader.vs", "world/glsl330/treeShader.fs")
	treeShader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(treeShader, "mvp"))
	treeShader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(treeShader, "instanceTransform"))
	treeModel.Materials.Shader = treeShader
	// treeModel.Transform = rl.MatrixScale(worldOption.tileSize, worldOption.tileSize, worldOption.tileSize)
	treeTranslations = make([]rl.Matrix, treeInstances)
	tmpIndex := 0
	for y := range world.tiles {
		for x := range world.tiles[y] {
			if world.tiles[y][x].tile == tree {
				xPos := (float32(x) - (width / 2)) * worldOption.tileSpacing
				yPos := (float32(y) - (length / 2)) * worldOption.tileSpacing
				height := float32(world.tiles[y][x].noise)*worldOption.tileSize - 1
				pos := rl.NewVector3(xPos, height*worldOption.heightMulti, yPos)
				mat := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
				treeTranslations[tmpIndex] = mat
				treeTranslations[tmpIndex] = rl.MatrixMultiply(treeTranslations[tmpIndex], rl.MatrixTranslate(0, 0, 0))
				tmpIndex++
			}
		}
	}

}

func (world *worldStruct) renderMesh() {
	rl.DrawMeshInstanced(*tileModel.Meshes, *tileModel.Materials, tileTranslations, tileInstances)
	for _, u := range treeTranslations {
		// log.Println(rl.Vector3Transform(rl.NewVector3(0, 0, 0), u))
		rl.DrawModel(treeModel, rl.Vector3Transform(rl.NewVector3(0, 0, 0), u), worldOption.tileSize, rl.White)
	}
	rl.DrawMeshInstanced(*treeModel.Meshes, *treeModel.Materials, treeTranslations, treeInstances)
}

func debugging() {
	rl.DrawText("Free camera default controls:", 20, 20, 10, rl.Black)
	rl.DrawText("- Mouse Wheel to Zoom in-out", 40, 40, 10, rl.DarkGray)
	rl.DrawText("- Mouse Wheel Pressed to Pan", 40, 60, 10, rl.DarkGray)
	rl.DrawText("- Alt + Mouse Wheel Pressed to Rotate", 40, 80, 10, rl.DarkGray)
	rl.DrawText("- Alt + Ctrl + Mouse Wheel Pressed for Smooth Zoom", 40, 100, 10, rl.DarkGray)
	rl.DrawText("- Z to y += 10", 40, 120, 10, rl.DarkGray)
	// Debugging
	start := 250
	posX := int32(10)
	posY := int32(10)
	incrY := func() int {
		start += 10
		return start
	}
	rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), posX, int32(incrY()), posY, rl.Black)
	rl.DrawText(fmt.Sprintf("Pos: %v", cam.Position), posX, int32(incrY()), posY, rl.Black)

}
