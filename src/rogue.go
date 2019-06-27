package main

import (
	blt "bearlibterminal"
	"strconv"
	"entity"
	"gamemap"
	"fov"
	"math/rand"
	"dijkstramaps"
	"fmt"
	"runtime"
)

const (
	WindowSizeX = 30
	WindowSizeY = 20
	MapWidth = WindowSizeX
	MapHeight = WindowSizeY
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	//FontSize = 16
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
	gameMapSrc = [][]int{}
	fieldOfView *fov.FieldOfVision
	dmap *dijkstramaps.EntityDijkstraMap
)

func init() {
	blt.Open()

	// BearLibTerminal uses configuration strings to set itself up, so we need to build these strings here
	// First set up the string for window properties (size and title)
	size := "size=120x46"
	title := "title='" + Title + "'"
	window := "window: " + size + "," + title

	// Next, setup the font config string
	//fontSize := "size=" + strconv.Itoa(FontSize)
	//font := "font: " + Font + ", " + fontSize
	font := "font: fonts/UbuntuMono.ttf, size=8x16; 0x1000: fonts/Floor.png, size=32x32, align=top-left; 0x2000: fonts/Humans.png, size=32x32, align=top-left;0x3001: fonts/HP.png, size=32x3, align=top-left"

	// Now, put it all together
	blt.Set(window + "; " + font)
	blt.Clear()

	// Create a GameMap, and initialize it
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMapSrc := make([][]int, MapHeight)
  for i := 0; i < MapHeight; i++ {
    gameMapSrc[i] = make([]int, MapWidth)
  }
	gameMap.InitializeMap()
	gameMap.GenerateRooms(gameMapSrc)
	gameMap.GenerateArena(gameMapSrc)
	
	player = &entity.GameEntity{X: 3, Y: 3, Layer: 3, Char: 0x2001, Color: "white", NPC: false, Name: "Player", HP: []int{40,40,}}
	npc := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "white", NPC: true, Name: "NPC 1", HP: []int{20,20,}}
	npc2 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "red", NPC: true, Name: "NPC 2", HP: []int{20,20,}}
	npc3 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "blue", NPC: true, Name: "NPC 3", HP: []int{20,20,}}
	entities = append(entities, player, npc, npc2, npc3)
  for _, e := range entities {
		e.X = rand.Intn(MapWidth - 2) + 1
		e.Y = rand.Intn(MapHeight - 2) + 1
		for ok := true; ok; ok = gameMap.IsBlocked(e.X, e.Y) {
			e.X = rand.Intn(MapWidth - 2) + 1
			e.Y = rand.Intn(MapHeight - 2) + 1
  	}
  }
	
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(5)
	
	dmap = dijkstramaps.NewEntityMap(3, player.X, player.Y, "test", MapWidth, MapHeight)
	dmap.GenerateMap(gameMap)
}
	
func main() {
	// Main game loop
	renderAll()
	defer recovery()
	runtime.LockOSThread()
	for {
    blt.Refresh()
    blt.Clear()
		key := blt.Read()

		// Clear each Entity off the screen
		for _, e := range entities {
			e.Clear()
		}

		if key != blt.TK_CLOSE {
			handleInput(key, player)
		} else {
			break
		}
		renderAll()
	}

	blt.Close()
}

func handleInput(key int, player *entity.GameEntity) {
	// Handle basic character movement in the four main directions

	var (
		dx, dy int
	)

	switch key {
	case blt.TK_RIGHT , blt.TK_KP_6:
		dx, dy = 1, 0
	case blt.TK_LEFT , blt.TK_KP_4:
		dx, dy = -1, 0
	case blt.TK_UP , blt.TK_KP_8:
		dx, dy = 0, -1
	case blt.TK_DOWN , blt.TK_KP_2:
		dx, dy = 0, 1
	case blt.TK_KP_9:
		dx, dy = 1, -1
	case blt.TK_KP_7:
		dx, dy = -1, -1
	case blt.TK_KP_1:
		dx, dy = -1, 1
	case blt.TK_KP_3:
		dx, dy = 1, 1
	}
	// Check to ensure that the tile the player is trying to move in to is a valid move (not blocked)
	if !gameMap.IsBlocked(player.X + dx, player.Y + dy) && !gameMap.Tiles[player.X + dx][player.Y + dy].Mob {
		player.Move(dx, dy)
		dmap.UpdateSourceCoordinates(player.X, player.Y)
		dmap.UpdateMap(gameMap)
	}
	IsMob(dx,dy)
}

func renderEntities() {
	// Draw every Entity present in the game. This gets called on each iteration of the game loop.
	for _, e := range entities {
		if e.NPC{ gameMap.Tiles[e.X][e.Y].Mob = false
		}
		e.Hunting(dmap, gameMap, player)
		if e.NPC{ gameMap.Tiles[e.X][e.Y].Mob = true
		}
		if gameMap.Tiles[e.X][e.Y].Visible {
	    e.Draw()
			blt.Layer(1)
			blt.Color(blt.ColorFromName(gameMap.Tiles[e.X][e.Y].Color))
			blt.Put(e.X*4, e.Y*2, gameMap.Tiles[e.X][e.Y].Symbol)
    }
		dmap.UpdateMap(gameMap)
	}
	blt.Print(1, 42, "Player X: " + strconv.Itoa(player.X) + " Y: " + strconv.Itoa(player.Y) + " HP:" + strconv.Itoa(player.HP[0]) + "/" + strconv.Itoa(player.HP[1]))
//	blt.Print(1, 44, "npc dmap: " + strconv.Itoa(dmap.ValuesMap[npc.X][npc.Y]))
//	blt.Print(1, 45, "npc2 dmap: " + strconv.Itoa(dmap.ValuesMap[npc2.X][npc2.Y]))
//	blt.Print(1, 46, "npc3 dmap: " + strconv.Itoa(dmap.ValuesMap[npc3.X][npc3.Y]))
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	fieldOfView.RayCast(player.X, player.Y, gameMap)

	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Color(blt.ColorFromName(gameMap.Tiles[x][y].Color))
      blt.Put(x*4, y*2, gameMap.Tiles[x][y].Symbol)
		}
	}
}

func renderDmap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Print(x + 40, y + 46, strconv.Itoa(dmap.ValuesMap[x][y]))
		}
	}
}

func renderBlock() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			if gameMap.Tiles[x][y].Mob{ 
				blt.Print(x + 20, y + 42, "1")
			} else { 
				blt.Print(x + 20, y + 42, "0")
			}
		}
	}
}

func renderAll() {
	// Convenience function to render all entities, followed by rendering the game map
	renderMap()
	renderEntities()
//	renderDmap()
}

func recovery() {
    if r := recover(); r != nil {
        fmt.Println("panic occured: ", r)
    }

}

func IsMob(dx , dy int) {
	var entitiestemp []*entity.GameEntity
	entitiestemp = append(entitiestemp, player)
	for _, e := range entities {
		if e.X == (player.X+dx) && e.Y ==(player.Y+dy) && e.NPC {
			blt.Print(1,44,player.Fight(gameMap, e))
		}
  	if e.HP[0] > 0{
		  entitiestemp = append(entitiestemp, e)
  	}
	}
	entities = entitiestemp
}
