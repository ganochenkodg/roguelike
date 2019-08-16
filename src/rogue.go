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
	"namegen"
	"time"
)

const (
	WindowSizeX = 23
	WindowSizeY = 15
	Title = "GoRogue"
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
	gameMapSrc = [][]int{}
	fieldOfView *fov.FieldOfVision
	dmap *dijkstramaps.EntityDijkstraMap
	gameMessages = messages.Messages{" "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," "," ",}
	gameState = "menu"
	MapWidth = WindowSizeX
	MapHeight = WindowSizeY
	
)

func init() {
	blt.Open()
	size := "size="+strconv.Itoa(WindowSizeX * 4)+"x"+strconv.Itoa(WindowSizeY * 2 + 6)
	title := "title='" + Title + "'"
	window := "window: " + size + "," + title
	font := "font: fonts/UbuntuMono.ttf, size=8x16; 0x1000: fonts/Dungeon.png, size=32x32, align=top-left; 0x2000: fonts/monsters.png, size=32x32, align=top-left;0x3001: fonts/HP.png, size=32x3, align=top-left"

	blt.Set(window + "; " + font)
	blt.Clear()

	//NewGame()
}
func DrawMenu() {
	blt.Color(blt.ColorFromName("white"))
	blt.Layer(1)
	for x := 5; x < 18; x++ {
		for y := 5; y < 14; y++ {
			blt.Put(x*4, y*2, 0x1000)
		}
	}
	blt.Color(blt.ColorFromName("#606060"))
	for x := 6; x < 17; x++ {
		for y := 6; y < 13; y++ {
			blt.Put(x*4, y*2, 0x1000)
		}
	}
	if gameState == "menu" {
		blt.Layer(2)
		blt.Color(blt.ColorFromName("white"))
		blt.Print(28, 14, "GO Roguelike v 0.01")
		blt.Print(28, 16, "Press any key to start New Game.")
		blt.Print(28, 18, "Press ESC to exit.")
	}
}

func NewGame(){
	rand.Seed( time.Now().UTC().UnixNano())
	MapWidth = rand.Intn(40) + 50
	MapHeight = rand.Intn(25) + 20
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMapSrc := make([][]int, MapHeight)
	for i := 0; i < MapHeight; i++ {
		gameMapSrc[i] = make([]int, MapWidth)
	}
	gameMap.InitializeMap()
	gameMap.GenerateRooms(gameMapSrc)
	gameMap.GenerateArena(gameMapSrc)
	//делаем игрока, 3 моба и раскидываем всех по карте на незанятые точки
	player = &entity.GameEntity{X: 3, Y: 3, Layer: 3, Char: 0x2000, Color: "white", NPC: false, Name: "Player", HP: []int{40,40,}, Vision: 5, Speed: 10}
	npc := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2001, Color: "white", NPC: true, Name: "NPC 1", HP: []int{20,20,}, Vision: 9, Speed: 10, SpeedPool: 0.0}
	npc2 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2003, Color: "red", NPC: true, Name: "NPC 2", HP: []int{20,20,}, Vision: 5, Speed: 6, SpeedPool: 0.0}
	npc3 := &entity.GameEntity{X: 28, Y: 5, Layer: 2, Char: 0x2004, Color: "blue", NPC: true, Name: "NPC 3", HP: []int{20,20,}, Vision: 15, Speed: 20, SpeedPool: 0.0}
	player.Name = namegen.GenerateName()
	newentities := append(entities, player, npc, npc2, npc3)
	entities = newentities
	for _, e := range entities {
		e.X = rand.Intn(MapWidth - 2) + 1
		e.Y = rand.Intn(MapHeight - 2) + 1
		for ok := true; ok; ok = gameMap.Tiles[e.X][e.Y].IsBlock() {
			e.X = rand.Intn(MapWidth - 2) + 1
			e.Y = rand.Intn(MapHeight - 2) + 1
		}
	}
	gameMap.CameraX, gameMap.CameraY = player.X, player.Y
	
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(player.Vision)
	
	dmap = dijkstramaps.NewEntityMap(3, player.X, player.Y, "test", MapWidth, MapHeight)
	dmap.GenerateMap(gameMap)
}
	
func main() {
	// игра периодически паникует, затычки от фаталов ниже. с ними норм
	//renderAll()
	DrawMenu()
	defer recovery()
	runtime.LockOSThread()
	for {
    blt.Refresh()
    blt.Clear()
		key := blt.Read()


		if key != blt.TK_CLOSE {
			handleInput(key, player)
		} else {
			break
		}
	}

	blt.Close()
}
//обрабокта нажатия клавиатуры в каждом цикле
func handleInput(key int, player *entity.GameEntity) {

	var (
		dx, dy int
	)
	if gameState == "menu" {
	DrawMenu()
	switch key {
	case blt.TK_ESCAPE:
		blt.Close()
	default:
		gameState = "game"
		blt.Clear()
		NewGame()
		renderAll()
	}
	return
	}
	
	if gameState == "journal" {
	gameMessages.DrawJournal()
	switch key {
	case blt.TK_ESCAPE, blt.TK_J:
		gameState = "game"
		blt.Clear()
		clearMobs()
		renderAll()
		return
	}
  return
	}
	if gameState == "map" {
	renderMapScreen()
	switch key {
	case blt.TK_ESCAPE, blt.TK_M:
		gameState = "game"
		blt.Clear()
		clearMobs()
		renderAll()
		return
	}
	return
	}
	
	if gameState == "game" {
	// каждый ход стираем игрока и мобов, потмо в новом месте нарисуютс как походят
	clearMobs()
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
	case blt.TK_J:
		gameState = "journal"
		gameMessages.DrawJournal()
		return
	case blt.TK_M:
		gameState = "map"
		renderMapScreen()
		return
	}
	
	// проверяем что там куда шагаем не стена и не нпс
	if !gameMap.Tiles[player.X + dx][player.Y + dy].IsBlock() && !gameMap.Tiles[player.X + dx][player.Y + dy].Mob {
		player.Move(dx, dy)
		gameMap.CameraX, gameMap.CameraY = player.X, player.Y
		dmap.UpdateSourceCoordinates(player.X, player.Y)
		dmap.UpdateMap(gameMap)
	}else if gameMap.Tiles[player.X + dx][player.Y + dy].IsDoor(){
		if rand.Intn(10) > 6{
			gameMap.Tiles[player.X + dx][player.Y + dy] = &gamemap.Tile{false, false, false, false, false,"white", 0x1004, player.X + dx, player.Y + dy}
		  dmap.UpdateMap(gameMap)
		  gameMessages.AddMessage("You open the door")
		} else {
			gameMessages.AddMessage("You cant open the door, try again")
		}
	}
	//если моб в направлении шага то драться
	IsMob(dx,dy)
	renderAll()
	return
  }
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
	    e.Draw(gameMap.CameraX, gameMap.CameraY)
			newx, newy := entity.GetCamera(e.X, e.Y, gameMap.CameraX, gameMap.CameraY)
			blt.Layer(1)
			blt.Color(blt.ColorFromName(gameMap.Tiles[e.X][e.Y].Color))
			blt.Put(newx*4, newy*2, gameMap.Tiles[e.X][e.Y].Symbol)
    }
		dmap.UpdateMap(gameMap)
	}
	blt.Print(50, 31, player.Name + " X: " + strconv.Itoa(player.X) + " Y: " + strconv.Itoa(player.Y) + " HP:" + strconv.Itoa(player.HP[0]) + "/" + strconv.Itoa(player.HP[1]))
}

func renderMap() {
	fieldOfView.RayCast(player.X, player.Y, gameMap)

	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Color(blt.ColorFromName(gameMap.Tiles[x][y].Color))
			newx, newy := entity.GetCamera(x, y, gameMap.CameraX, gameMap.CameraY)
      blt.Put(newx*4, newy*2, gameMap.Tiles[x][y].Symbol)
		}
	}
}

func renderMapScreen() {
	blt.Clear()
	blt.Layer(5)
	xoffset := (92 - gameMap.Width) / 2
	yoffset := (36 - gameMap.Height) / 2
	if xoffset + gameMap.CameraX > 91{
		xoffset = 91 - gameMap.CameraX
	}
	if xoffset + gameMap.CameraX < 1{
		xoffset = 1 - gameMap.CameraX
	}
	if yoffset + gameMap.CameraY > 33{
		yoffset = 33 - gameMap.CameraY
	}
	if yoffset + gameMap.CameraY < 2{
		yoffset = 2 - gameMap.CameraY
	}
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			blt.Color(blt.ColorFromName(gameMap.Tiles[x][y].Color))
      switch gameMap.Tiles[x][y].Symbol {
			case 0x1001 , 0x1002:
				blt.Print(x + xoffset, y + yoffset, ".")
			case 0x1003:
				blt.Print(x + xoffset, y + yoffset, "|")
			case 0x1004:
				blt.Print(x + xoffset, y + yoffset, "+")
			case 0x1000:
				blt.Print(x + xoffset, y + yoffset, "#")
			}
		}
	}
	blt.Color(blt.ColorFromName("white"))
	blt.ClearArea(0,0,92,2)
	blt.ClearArea(0,34,92,2)
	blt.Print(41,1,"Game map:")
	blt.Print(gameMap.CameraX + xoffset, gameMap.CameraY + yoffset, "@")
  blt.Print(1,34,"Press M or ESC to return")
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
func clearMobs() {
	for _, e := range entities {
	 e.Clear(gameMap.CameraX, gameMap.CameraY)
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
