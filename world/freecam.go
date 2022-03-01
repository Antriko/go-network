package world

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var world [][]mapTile
var cam rl.Camera
var mapMatrix rl.Matrix

type worldOptionStruct struct {
	tileHeight  float32
	tileSpacing float32
}

var worldOption worldOptionStruct

func Freecam() {
	rl.InitWindow(1280, 720, "Freecam")
	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	world = createWorld(5)
	worldOption = worldOptionStruct{1, 1}
	meshGen()

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

		rl.DrawGrid(20, 5)
		// renderWorld(world)
		renderMesh(world)

		rl.EndMode3D()
		debugging()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func renderWorld(world [][]mapTile) {
	spacing := float32(5.0)
	width := float32(len(world))
	length := float32(len(world[0]))
	for y := range world {
		for x := range world[y] {

			xPos := (float32(x) - (width / 2)) * worldOption.tileSpacing
			yPos := (float32(y) - (length / 2)) * worldOption.tileSpacing

			var col color.RGBA
			// if (x+y)%2 == 0 {
			// 	col = rl.NewColor(2, 35, 28, 255)
			// } else {
			// 	col = rl.NewColor(0, 77, 37, 255)
			// }

			if world[y][x].noise < 0.15 {
				col = rl.NewColor(0, 45, 57, 255)
			} else if world[y][x].noise < 0.3 {
				col = rl.NewColor(0, 163, 204, 255)
			} else if world[y][x].noise < 0.6 {
				col = rl.NewColor(95, 99, 68, 255)
			} else {
				col = rl.NewColor(65, 69, 47, 255)
			}

			height := float32(world[y][x].noise) * worldOption.tileHeight
			if height < 0.3*worldOption.tileHeight {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height, yPos)
			rl.DrawCube(pos, spacing, worldOption.tileHeight, spacing, col)
			rl.DrawCubeWires(pos, spacing, worldOption.tileHeight, spacing, rl.Black)
		}
	}
}

var instances int
var translations []rl.Matrix // Locations of instances
var tileMesh rl.Model

func meshGen() {

	// Collection of meshes
	// All are combined into one model
	// Shader returns the colour of mesh depending on height level

	instances = len(world) * len(world)
	translations = make([]rl.Matrix, instances) // Locations of instances;

	// Create basic cube mesh
	tileMesh = rl.LoadModelFromMesh(rl.GenMeshCube(worldOption.tileSpacing, worldOption.tileHeight, worldOption.tileSpacing))
	tex := rl.GenImageChecked(2, 2, 1, 1, rl.Red, rl.Green)
	texture := rl.LoadTextureFromImage(tex)
	rl.SetMaterialTexture(tileMesh.Materials, rl.MapDiffuse, texture)
	rl.UnloadImage(tex)

	// Shader
	shader := rl.LoadShader("world/glsl330/basic.vs", "world/glsl330/basic.fs")
	shader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))
	tileMesh.Materials.Shader = shader
	// rl.transform

	// Assign location for all cubes to create map
	width := float32(len(world))
	length := float32(len(world[0]))
	for y := range world {
		for x := range world[y] {
			xPos := (float32(x) - (width / 2)) * worldOption.tileSpacing
			yPos := (float32(y) - (length / 2)) * worldOption.tileSpacing
			height := float32(world[y][x].noise) * worldOption.tileHeight
			if height < 0.3*worldOption.tileHeight {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height, yPos)

			index := y*int(width) + x
			mat := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
			translations[index] = mat
			translations[index] = rl.MatrixMultiply(translations[index], rl.MatrixTranslate(0, 0, 0))

			// log.Println(index, instances, rl.Vector3Transform(rl.NewVector3(0, 0, 0), translations[index]).Y)
		}
	}
}

func renderMesh(world [][]mapTile) {

	// rl.DrawMesh(*tileMesh.Meshes, *tileMesh.Materials, rl.MatrixIdentity())
	rl.DrawMeshInstanced(*tileMesh.Meshes, *tileMesh.Materials, translations, instances)

	// for _, u := range translations {
	// 	rl.DrawMesh(*tileMesh.Meshes, *tileMesh.Materials, u)
	// 	// rl.DrawModel(tileMesh, rl.Vector3Transform(rl.NewVector3(0, 0, 0), u), 1, rl.White)
	// }
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
