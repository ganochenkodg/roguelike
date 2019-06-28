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
	"messages"
)

const (
	WindowSizeX = 30
	WindowSizeY = 20
	MapWidth = WindowSizeX
	MapHeight = WindowSizeY
	Title = "GoRogue"
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
	gameMapSrc = [][]int{}
	fieldOfView *fov.FieldOfVision
	dmap *dijkstramaps.EntityDijkstraMap
	gameMessages = messages.Messages{" "," "," "," "," ",}
	
)

func init() {
	blt.Open()

	size := "size="+strconv.Itoa(WindowSizeX * 4)+"x"+strconv.Itoa(WindowSizeY * 2 + 6)
	title := "title='" + Title + "'"
	window := "window: " + size + "," + title
	font := "font: fonts/UbuntuMono.ttf, size=8x16; 0x1000: fonts/Floor.png, size=32x32, align=top-left; 0x2000: fonts/Humans.png, size=32x32, align=top-left;0x3001: fonts/HP.png, size=32x3, align=top-left"

	blt.Set(window + "; " + font)
	blt.Clear()

	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMapSrc := make([][]int, MapHeight)
  for i := 0; i < MapHeight; i++ {
    gameMapSrc[i] = make([]int, MapWidth)
  }
	gameMap.InitializeMap()
	gameMap.GenerateRooms(gameMapSrc)
	gameMap.GenerateArena(gameMapSrc)
	//делаем игрока, 3 моба и раскидываем всех по карте на незанятые точки
	player = &entity.GameEntity{X: 3, Y: 3, Layer: 3, Char: 0x2001, Color: "white", NPC: false, Name: "Player", HP: []int{40,40,}, Vision: 5}
	npc := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "white", NPC: true, Name: "NPC 1", HP: []int{20,20,}, Vision: 9}
	npc2 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "red", NPC: true, Name: "NPC 2", HP: []int{20,20,}, Vision: 5}
	npc3 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2002, Color: "blue", NPC: true, Name: "NPC 3", HP: []int{20,20,}, Vision: 15}
	entities = append(entities, player, npc, npc2, npc3)
  for _, e := range entities {
		e.X = rand.Intn(MapWidth - 2) + 1
		e.Y = rand.Intn(MapHeight - 2) + 1
		for ok := true; ok; ok = gameMap.Tiles[e.X][e.Y].IsBlock() {
			e.X = rand.Intn(MapWidth - 2) + 1
			e.Y = rand.Intn(MapHeight - 2) + 1
  	}
  }
	
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(player.Vision)
	
	dmap = dijkstramaps.NewEntityMap(3, player.X, player.Y, "test", MapWidth, MapHeight)
	dmap.GenerateMap(gameMap)
	
}
	
func main() {
	// игра периодически паникует, затычки от фаталов ниже. с ними норм
	renderAll()
	defer recovery()
	runtime.LockOSThread()
	for {
    blt.Refresh()
    blt.Clear()
		key := blt.Read()

		// каждый ход стираем игрока и мобов, потмо в новом месте нарисуютс как походят
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
//обрабокта нажатия клавиатуры в каждом цикле
func handleInput(key int, player *entity.GameEntity) {

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
	// проверяем что там куда шагаем не стена и не нпс
	if !gameMap.Tiles[player.X + dx][player.Y + dy].IsBlock() && !gameMap.Tiles[player.X + dx][player.Y + dy].Mob {
		player.Move(dx, dy)
		dmap.UpdateSourceCoordinates(player.X, player.Y)
		dmap.UpdateMap(gameMap)
	}
	//если моб в направлении шага то драться
	IsMob(dx,dy)
}

func renderEntities() {
	for _, e := range entities {
		if e.NPC{ gameMap.Tiles[e.X][e.Y].Mob = false
		}
		e.Hunting(dmap, gameMap, player, gameMessages)
		//казалось бы нафига подряд 2 ифа с одинаковым условием. если после битвы нпс жив то клетку снвоа проставляем как Mob
		if e.NPC{ gameMap.Tiles[e.X][e.Y].Mob = true
		}
		//рисуем если игрок видит их
		if gameMap.Tiles[e.X][e.Y].Visible {
	    e.Draw()
			blt.Layer(1)
			blt.Color(blt.ColorFromName(gameMap.Tiles[e.X][e.Y].Color))
			blt.Put(e.X*4, e.Y*2, gameMap.Tiles[e.X][e.Y].Symbol)
    }
		dmap.UpdateMap(gameMap)
	}
	blt.Print(60, 41, "Player X: " + strconv.Itoa(player.X) + " Y: " + strconv.Itoa(player.Y) + " HP:" + strconv.Itoa(player.HP[0]) + "/" + strconv.Itoa(player.HP[1]))
}

func renderMap() {
	fieldOfView.RayCast(player.X, player.Y, gameMap)

	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Color(blt.ColorFromName(gameMap.Tiles[x][y].Color))
      blt.Put(x*4, y*2, gameMap.Tiles[x][y].Symbol)
		}
	}
}

//отладочная функция, нарисовать карту путей
func renderDmap() {
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Print(x + 40, y + 46, strconv.Itoa(dmap.ValuesMap[x][y]))
		}
	}
}

//отладочная функция, нарисовать карту стен
func renderBlock() {
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
	renderMap()
	renderEntities()
	gameMessages.PrintMessages()
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
			gameMessages.AddMessage(player.Fight(gameMap, e))
		}
  	if e.HP[0] > 0{
		  entitiestemp = append(entitiestemp, e)
  	}
	}
	entities = entitiestemp
}
