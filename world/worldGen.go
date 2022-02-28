package world

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gookit/color"
	noise "github.com/ojrac/opensimplex-go"
)

func createWorld(size int) [][]mapTile {
	log.SetFlags(log.Lshortfile)
	log.Println(mapGreen(" - WORLD GEN - "))

	// Create empty world
	world := worldGen(size)
	// Add objects
	populateWorld(world)

	return world
}

type tileType uint32

const (
	empty tileType = iota
	tree
)

type mapTile struct {
	noise float64
	tile  tileType
}

var mapBlue = color.New(color.FgBlack, color.BgLightBlue).Render
var mapGreen = color.New(color.FgBlack, color.BgLightGreen).Render
var mapDarkGreen = color.New(color.FgBlack, color.BgHiGreen).Render
var mapRed = color.New(color.FgBlack, color.BgRed).Render

// Keep size odd so gradients work better
// personal note; 91 max for own console width
func worldGen(size int) [][]mapTile {

	width, height := size, size
	rand.Seed(time.Now().UnixNano())

	// Create basic simplex noise
	world := combineNoise(width, height, combineNoise(width, height, noiseGen(width, height), noiseGen(width, height)), combineNoise(width, height, noiseGen(width, height), noiseGen(width, height)))

	// 60% at least >0.3
	if !checkValid(size, 0.3, 0.6, world) {
		// recursion, attempt to create new world
		// printNoise(world)
		log.Println(mapRed(" WORLD NOT VALID - GEN NEW WORLD "))
		world = worldGen(size)
	}
	return world
}

func printNoise(noise [][]mapTile) {
	fmt.Println()
	for y := range noise {
		for x := range noise[y] {
			hex := strconv.FormatInt(int64(math.Round(noise[y][x].noise*10)), 16)

			if noise[y][x].noise < 0.3 {
				fmt.Printf(mapBlue(" %v "), hex)
			} else if noise[y][x].noise < 0.6 {
				// tmp := math.Round(noise[y][x] * 100)
				fmt.Printf(mapGreen(" %v "), hex)
			} else {
				fmt.Printf(mapDarkGreen(" %v "), hex)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func printEmpty(noise [][]mapTile) {
	fmt.Println()
	for y := range noise {
		for x := range noise[y] {
			if noise[y][x].noise < 0.3 {
				fmt.Printf(mapBlue("   "))
			} else if noise[y][x].noise < 0.6 {
				// tmp := math.Round(noise[y][x] * 100)
				fmt.Printf(mapGreen("   "))
			} else {
				fmt.Printf(mapDarkGreen("   "))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func printTile(tile [][]mapTile) {
	fmt.Println()
	for y := range tile {
		for x := range tile[y] {
			if tile[y][x].noise < 0.3 {
				fmt.Printf(mapBlue(" %d "), tile[y][x].tile)
			} else if tile[y][x].noise < 0.6 {
				// tmp := math.Round(noise[y][x] * 100)
				fmt.Printf(mapGreen(" %d "), tile[y][x].tile)
			} else {
				fmt.Printf(mapDarkGreen(" %d "), tile[y][x].tile)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func noiseGen(width int, height int) [][]mapTile {
	// Create basic simplex noise
	noiseWorld := make([][]mapTile, width)
	for y := range noiseWorld {
		noiseWorld[y] = make([]mapTile, height)
	}

	noiseInst := noise.New(rand.Int63()) // Produces smaller islands but is less successful in validation
	// noiseInst := noise.NewNormalized(rand.Int63()) // Islands are bigger and is more likely attached to edge

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			xFloat := float64(x) / float64(width)
			yFloat := float64(y) / float64(height)
			noiseWorld[y][x].noise = noiseInst.Eval2(xFloat, yFloat)
		}
	}
	return noiseWorld
}

func combineNoise(width int, height int, noise1 [][]mapTile, noise2 [][]mapTile) [][]mapTile {
	// check if width/height is same

	for y := range noise1 {
		for x := range noise1[y] {
			// Make combination not as strong
			grad := noise1[y][x].noise - (noise2[y][x].noise / 1.5)
			if grad < 0 {
				grad = 0
			} else if grad > 1 {
				grad = 1
			}
			noise1[y][x].noise = grad
		}
	}
	return noise1
}

func mergeCircleNoise(width int, height int, noise [][]mapTile) [][]mapTile {
	// Create circle gradient to factor out edges to create island
	circle := make([][]mapTile, width)
	for y := range noise {
		circle[y] = make([]mapTile, height)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Sin() depending on axis of XY
			xAxis := float64((180 / (width - 1)) * x)
			xGrad := 1 - math.Sin(xAxis*(math.Pi/180))

			yAxis := float64((180 / (height - 1)) * y)
			yGrad := 1 - math.Sin(yAxis*(math.Pi/180))

			// Add them both together to get a circle gradient
			div := 2.0 // default 2; more = less strong
			circle[y][x].noise = (xGrad / div) + (yGrad / div)
		}
	}

	// printNoise(circle)
	return combineNoise(width, height, noise, circle)
}
func checkValid(size int, threshold float64, percent float32, noise [][]mapTile) bool {
	// check if total noise contains >= percent%
	// dont want mainly water generated islands
	total := 0.0
	for y := range noise {
		for x := range noise[y] {
			if noise[y][x].noise > threshold {
				total++
			}
		}
	}
	// log.Println(total, size*size, float64(size*size)*float64(percent))
	if total <= (float64(size*size))*float64(percent) {
		return false
	}
	return true
}

func populateWorld(world [][]mapTile) {
	populateTrees(world)
}

func populateTrees(world [][]mapTile) {

}
