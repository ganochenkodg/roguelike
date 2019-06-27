package fov

import (
	"math"
	"gamemap"
	//"fmt"
	//"strconv"
)

type FieldOfVision struct {
	cosTable map[int]float64
	sinTable map[int]float64
	torchRadius int
}

func (f *FieldOfVision) Initialize() {

	f.cosTable = make(map[int]float64)
	f.sinTable = make(map[int]float64)

	for i := 0; i < 360; i++ {
		ax := math.Sin(float64(i) / (float64(180) / math.Pi))
		ay := math.Cos(float64(i) / (float64(180) / math.Pi))

		f.sinTable[i] = ax
		f.cosTable[i] = ay
	}
}

func (f *FieldOfVision) SetTorchRadius(radius int) {
	if radius > 0 {
		f.torchRadius = radius
	}
}

func (f *FieldOfVision) RayCast(playerX, playerY int, gameMap *gamemap.Map) {
	// Cast out rays each degree in a 360 circle from the player. If a ray passes over a floor (does not block sight)
	// tile, keep going, up to the maximum torch radius (view radius) of the player. If the ray intersects a wall
	// (blocks sight), stop, as the player will not be able to see past that. Every visible tile will get the Visible
	// and Explored properties set to true.	
	
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			gameMap.Tiles[x][y].Visible = false
			if gameMap.Tiles[x][y].Visited {
				gameMap.Tiles[x][y].Color = "#808080"
			} else {
				gameMap.Tiles[x][y].Color = "black"
			}
		}
	}	
	
	gameMap.Tiles[playerX][playerY].Visited = true
	gameMap.Tiles[playerX][playerY].Color = "white"

	for i := 0; i < 360; i ++ {

		ax := f.sinTable[i]
		ay := f.cosTable[i]

		x := float64(playerX)
		y := float64(playerY)

		// Mark the players current position as explored
		for j := 0; j < f.torchRadius; j++ {
			x -= ax
			y -= ay

			if x < 0 || x > float64(gameMap.Width - 1) || y < 0 || y > float64(gameMap.Height - 1) {
				// If the ray is cast outside of the map, stop
				break
			}

			gameMap.Tiles[int(x)][int(y)].Visited = true
			gameMap.Tiles[int(x)][int(y)].Visible = true
			switch f.torchRadius - j {
			case 1: gameMap.Tiles[int(x)][int(y)].Color = "#E8E8E8"
			case 2: gameMap.Tiles[int(x)][int(y)].Color = "#F7F7F7"
			default: gameMap.Tiles[int(x)][int(y)].Color = "#FFFFFF"
			}

			if gameMap.Tiles[int(x)][int(y)].Blocks_sight == true {
				// The ray hit a wall, go no further
				break
			}
		}
	}
}
