package server

import (
	"fmt"

	"github.com/gookit/color"
)

var blue = color.New(color.FgBlack, color.BgLightBlue).Render

func worldGen() {

	width, height := 50, 50

	world := make([][]int, width)
	for x := range world {
		world[x] = make([]int, height)
	}

	printWorld(world)
	/*
		World size 50x50
		Small island, outer edge covered by water

		Each tile contains {
			biome	: grid	-	forest/plains/desert/water
			object	: model	-	tree/rock/cactus/structure
		}
	*/
}

func printWorld(world [][]int) {
	for x := range world {
		for y := range world[x] {
			switch world[x][y] {
			case 0:
				fmt.Print(blue(" - "))
			}
		}
		fmt.Println()
	}
}
