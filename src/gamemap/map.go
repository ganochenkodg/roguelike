package gamemap

import (
	"math/rand"
	"time"
	"dungeon"
)

type Tile struct {
	Blocked bool
	Blocks_sight bool
	Visited bool
	Visible bool
	Mob bool
	Color string
	Symbol int
	X int
	Y int
}

type Map struct {
	Width  int
	Height int
	CameraX int
	CameraY int
	Tiles  [][]*Tile
}

func (m *Map) InitializeMap() {
	// генерим двухмерный массив нужных размеров и инициализируем генератор рандома
	m.Tiles = make([][]*Tile, m.Width)
	for i := range m.Tiles {
		m.Tiles[i] = make([]*Tile, m.Height)
	}
	rand.Seed( time.Now().UTC().UnixNano())
}

//заполняем карту согласно сгенернному подземелью. добавить больше типов клеток, туманы там вские...
func (m *Map) GenerateArena(src [][]int) {
	rand.Seed( time.Now().UTC().UnixNano())
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			switch src[y][x] {
			case 11:
				m.Tiles[x][y] = &Tile{true, true, false, false, false,"white", 0x1000, x, y}
			case 0:
				//иеогда рисуем полу другой спрайт, с трещинами и дырами
				if rand.Intn(100) < 98{ 
					m.Tiles[x][y] = &Tile{false, false, false, false, false,"white", 0x1001, x, y}
				} else {
					m.Tiles[x][y] = &Tile{false, false, false, false, false,"white", 0x1002, x, y}
				}
			}
		}
	}
	doorx , doory := 0, 0
	for x := 0; x < rand.Intn(5) + 10; x++{
	doorx = rand.Intn(m.Width - 2) + 1
	doory = rand.Intn(m.Height - 2) + 1
	for ok := true; ok; ok = m.IsDoorReady(doorx,doory) {
		doorx = rand.Intn(m.Width - 2) + 1
		doorx = rand.Intn(m.Height - 2) + 1
	}
	m.Tiles[doorx][doory] = &Tile{true, true, false, false, false,"white", 0x1003, doorx, doory}
  }
}

//черная магия
func (m *Map) GenerateRooms(src [][]int) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			src[y][x] = 11
		}
	}

	for x := 1; x < m.Width - 1; x++ {
		for y := 1; y < m.Height - 1; y++ {
			src[y][x] = 0
		}
	}
	rand.Seed( time.Now().UTC().UnixNano())
	dungeon := dungeon.NewDungeon(m.Width, m.Height, rand.Intn(20) + 10)
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if dungeon.Grid[x][y] == 0{
				src[y][x] = 11
			}
		}
	}

}

func (t *Tile) isVisited() bool {
	return t.Visited
}

func (t *Tile) IsWall() bool {
	if t.Blocks_sight && t.Blocked {
		return true
	} else {
		return false
	}
}

func (t *Tile) IsBlock() bool {
	if t.Blocked {
		return true
	} else {
		return false
	}
}

func (m *Map) IsDoorReady(x,y int) bool {
	if m.Tiles[x-1][y].Blocked && m.Tiles[x+1][y].Blocked && !(m.Tiles[x][y-1].Blocked) && !(m.Tiles[x][y+1].Blocked){
		return true
	} else if m.Tiles[x][y-1].Blocked && m.Tiles[x][y+1].Blocked && !(m.Tiles[x-1][y].Blocked) && !(m.Tiles[x+1][y].Blocked) {
		return true
	} else{
		return false
	}
}
