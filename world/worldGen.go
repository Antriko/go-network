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

func createWorld(size int) *worldStruct {
	log.SetFlags(log.Lshortfile)
	log.Println(mapGreen(" - WORLD GEN - "))

	// Create empty world
	world := &worldStruct{
		worldGen(size),
	}

	// Add objects
	world.populateWorld()

	return world
}

type tileType uint32

const thresholdWater = 0.3
const thresholdLand = 0.6

const (
	empty tileType = iota
	tree
)

type mapTile struct {
	noise float64
	tile  tileType
}
type worldStruct struct {
	tiles [][]mapTile
}

var mapBlue = color.New(color.FgBlack, color.BgLightBlue).Render
var mapGreen = color.New(color.FgBlack, color.BgLightGreen).Render
var mapDarkGreen = color.New(color.FgBlack, color.BgHiGreen).Render
var mapRed = color.New(color.FgBlack, color.BgRed).Render
var mapBrown = color.New(color.FgBlack, color.BgDarkGray).Render

// Keep size odd so gradients work better
// personal note; 91 max for own console width
func worldGen(size int) [][]mapTile {

	width, height := size, size
	rand.Seed(time.Now().UnixNano())

	// Create basic simplex noise
	worldGenerated := combineNoise(width, height, combineNoise(width, height, noiseGen(width, height), noiseGen(width, height)), combineNoise(width, height, noiseGen(width, height), noiseGen(width, height)))

	// 60% at least >0.3
	if !checkValid(size, 0.3, 0.6, worldGenerated) {
		// recursion, attempt to create new world
		log.Println(mapRed(" WORLD NOT VALID - GEN NEW WORLD "))
		worldGenerated = worldGen(size)
	}
	return worldGenerated
}

func printNoise(noise [][]mapTile) {
	fmt.Println()
	for y := range noise {
		for x := range noise[y] {
			hex := strconv.FormatInt(int64(math.Round(noise[y][x].noise*10)), 16)

			if noise[y][x].noise < thresholdWater {
				fmt.Printf(mapBlue(" %v "), hex)
			} else if noise[y][x].noise < thresholdLand {
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
			if noise[y][x].noise < thresholdWater {
				fmt.Printf(mapBlue("   "))
			} else if noise[y][x].noise < thresholdLand {
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

			if tile[y][x].tile == tree {
				fmt.Printf(mapBrown(" %d "), tile[y][x].tile)
			} else if tile[y][x].noise < thresholdWater {
				fmt.Printf(mapBlue(" %d "), tile[y][x].tile)
			} else if tile[y][x].noise < thresholdLand {
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
func (world *worldStruct) printWorld() {
	printTile(world.tiles)
}
func (world *worldStruct) printNoise() {
	printNoise(world.tiles)
}
func (world *worldStruct) printEmpty() {
	printEmpty(world.tiles)
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

func (world *worldStruct) populateWorld() {
	world.populateTrees()
}

func (world *worldStruct) populateTrees() {
	rand.Seed(time.Now().UnixNano())
	randomNoise := noiseGen(len(world.tiles), len(world.tiles))
	for y := range world.tiles {
		for x := range world.tiles[y] {
			if world.tiles[y][x].noise > thresholdWater {
				randomNum := rand.Intn(100)
				if randomNum <= 10 { // % threshold
					if randomNoise[y][x].noise > 0.1 {

						// check around to see if obstructed by water
						// O O O
						// O X O
						// O O O

						// 0 will be accounted for
						radius := 2
						diameter := (radius * 2)
						canSpawn := true
						for diaX := 0; diaX <= diameter; diaX++ {
							for diaY := 0; diaY <= diameter; diaY++ {

								if -radius+y+diaY < 0 || -radius+y+diaY > len(world.tiles)-1 || -radius+x+diaX < 0 || -radius+x+diaX > len(world.tiles)-1 {
									continue
								}
								if world.tiles[-radius+y+diaY][-radius+x+diaX].noise < thresholdWater {
									canSpawn = false
								}
							}
						}
						if canSpawn {
							world.tiles[y][x].tile = tree
						}
					}
				}
			}
		}
	}
	world.printWorld()
}
