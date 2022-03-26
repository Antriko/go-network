package world

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var world *WorldStruct
var cam rl.Camera

type WorldOptionStruct struct {
	TileSize    float32
	TileSpacing float32
	HeightMulti float32
}

var WorldOption WorldOptionStruct
var worldSize = 5

func Freecam() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "Freecam")
	rl.SetTraceLog(rl.LogError)
	rl.SetTargetFPS(60)
	rl.SetWindowPosition(3200, 100) // Stops displaying on my left monitor

	WorldOption = WorldOptionStruct{5, 5, 2}
	world = CreateWorld(worldSize)
	world.MeshGen(WorldOption)

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

		if rl.IsKeyPressed(rl.KeySpace) {
			world = CreateWorld(worldSize)
			world.MeshGen(WorldOption)
		} else if rl.IsKeyPressed(rl.KeyMinus) {
			worldSize--
			if worldSize < 2 {
				worldSize = 2
			}
			world = CreateWorld(worldSize)
			world.MeshGen(WorldOption)
		} else if rl.IsKeyPressed(rl.KeyEqual) {
			worldSize++
			world = CreateWorld(worldSize)
			world.MeshGen(WorldOption)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(cam)
		rl.DrawGrid(10, WorldOption.TileSpacing)
		world.RenderMesh()

		rl.EndMode3D()
		debugging()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

// Not efficiently rendered - deprecated
// Ignores shader rending
func (world *WorldStruct) renderWorld() {
	spacing := float32(5.0)
	width := float32(world.Size)
	length := width
	for y := range world.Tiles {
		for x := range world.Tiles[y] {

			xPos := (float32(x) - (width / 2)) * WorldOption.TileSpacing
			yPos := (float32(y) - (length / 2)) * WorldOption.TileSpacing

			var col color.RGBA
			// if (x+y)%2 == 0 {
			// 	col = rl.NewColor(2, 35, 28, 255)
			// } else {
			// 	col = rl.NewColor(0, 77, 37, 255)
			// }

			if world.Tiles[y][x].Noise < 0.15 {
				col = rl.NewColor(0, 45, 57, 255)
			} else if world.Tiles[y][x].Noise < 0.3 {
				col = rl.NewColor(0, 163, 204, 255)
			} else if world.Tiles[y][x].Noise < 0.6 {
				col = rl.NewColor(95, 99, 68, 255)
			} else {
				col = rl.NewColor(65, 69, 47, 255)
			}

			height := float32(world.Tiles[y][x].Noise) * WorldOption.TileSize
			if height < 0.3 {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height*WorldOption.TileSize, yPos)
			rl.DrawCube(pos, spacing, WorldOption.TileSize, spacing, col)
			rl.DrawCubeWires(pos, spacing, WorldOption.TileSize, spacing, rl.Black)
		}
	}
}

func (world *WorldStruct) MeshGen(WorldOption WorldOptionStruct) {
	treeInstances := 0
	// Collection of translations in which the GPU will render the same model over those positions
	// to save compute time from initialising the model over and over
	// Shader returns the colour of mesh depending on height level

	// Tile instances
	tileInst := Instance{}
	tileInst.Type = empty
	tileInst.Instances = world.Size * world.Size
	tileInst.Translations = make([]rl.Matrix, tileInst.Instances) // Locations of instances;
	// log.Println(WorldOption)
	tileInst.Model = rl.LoadModelFromMesh(rl.GenMeshCube(WorldOption.TileSpacing, WorldOption.TileSize, WorldOption.TileSpacing))

	// Shader
	tileShader := rl.LoadShader("world/glsl330/tileShader.vs", "world/glsl330/tileShader.fs")
	tileShader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(tileShader, "mvp"))
	tileShader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(tileShader, "instanceTransform"))
	tileInst.Model.Materials.Shader = tileShader

	// Assign location for all cubes to create map
	size := float32(world.Size)
	offset := float32(WorldOption.TileSize) / 2
	for y := range world.Tiles {
		for x := range world.Tiles[y] {
			xPos := (float32(x)-(size/2))*WorldOption.TileSpacing + offset
			yPos := (float32(y)-(size/2))*WorldOption.TileSpacing + offset
			height := float32(world.Tiles[y][x].Noise) * WorldOption.TileSize
			if height < 0.3*WorldOption.TileSize {
				height = 0.3
			}
			pos := rl.NewVector3(xPos, height*WorldOption.HeightMulti, yPos)

			index := y*int(size) + x
			mat := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
			tileInst.Translations[index] = mat
			tileInst.Translations[index] = rl.MatrixMultiply(tileInst.Translations[index], rl.MatrixTranslate(0, 0, 0))

			// log.Println(index, instances, rl.Vector3Transform(rl.NewVector3(0, 0, 0), translations[index]).Y)

			// Get total amount of trees to get amount of instances needed
			if world.Tiles[y][x].Tile == tree {
				treeInstances++
			}
		}
	}

	world.Instances[empty] = tileInst

	// Tree instances
	treeInst := Instance{}
	treeInst.Type = tree
	treeInst.Model = rl.LoadModel("models/world/Tree.glb")
	treeInst.Instances = treeInstances
	treeInst.Translations = make([]rl.Matrix, treeInst.Instances)
	tmpIndex := 0
	for y := range world.Tiles {
		for x := range world.Tiles[y] {
			if world.Tiles[y][x].Tile == tree {
				xPos := (float32(x)-(size/2))*WorldOption.TileSpacing + offset
				yPos := (float32(y)-(size/2))*WorldOption.TileSpacing + offset
				height := float32(world.Tiles[y][x].Noise) * WorldOption.TileSize
				pos := rl.NewVector3(xPos, height*WorldOption.HeightMulti, yPos)
				mat := rl.MatrixTranslate(pos.X, pos.Y, pos.Z)
				treeInst.Translations[tmpIndex] = mat
				treeInst.Translations[tmpIndex] = rl.MatrixMultiply(treeInst.Translations[tmpIndex], rl.MatrixTranslate(0, 0, 0))
				tmpIndex++
			}
		}
	}
	world.Instances[tree] = treeInst
}

func (world *WorldStruct) RenderMesh() {
	for _, instance := range world.Instances {
		// log.Println(instance.Type, len(instance.Translations), instance.Instances)
		if instance.Type == tree {
			for _, u := range instance.Translations {
				rl.DrawModel(instance.Model, rl.Vector3Transform(rl.NewVector3(0, 0, 0), u), WorldOption.TileSize*0.75, rl.White)
			}
			continue
		}
		rl.DrawMeshInstanced(*instance.Model.Meshes, *instance.Model.Materials, instance.Translations, instance.Instances)
	}

	// rl.DrawMeshInstanced(*tileModel.Meshes, *tileModel.Materials, tileTranslations, tileInstances)
	// for _, u := range world.Instances[tree] {
	// 	// log.Println(rl.Vector3Transform(rl.NewVector3(0, 0, 0), u))
	// 	rl.DrawModel(treeModel, rl.Vector3Transform(rl.NewVector3(0, 0, 0), u), WorldOption.TileSize*0.75, rl.White)
	// }
	// rl.DrawMeshInstanced(*treeModel.Meshes, *treeModel.Materials, treeTranslations, treeInstances)

	// rl.DrawModel(world.Instances[tree].Model, rl.NewVector3(0, 0, 0), 1, rl.White)
	// rl.DrawModel(world.Instances[empty].Model, rl.NewVector3(1, 1, 0), 1, rl.White)

	// Middle of map -
	// rl.DrawCube(rl.NewVector3(0, 0, 0), 1, 100, 1, rl.Red)
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
