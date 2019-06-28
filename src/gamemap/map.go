package gamemap

import (
	"math/rand"
	"time"
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
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			switch src[y][x] {
			case 12:
				m.Tiles[x][y] = &Tile{true, true, false, false, false, "black", 0, x, y}
			case 11:
				m.Tiles[x][y] = &Tile{true, true, false, false, false,"white", 0x1017, x, y}
			case 0:
				//иеогда рисуем полу другой спрайт, с трещинами и дырами
				if rand.Intn(100) < 98{ 
					m.Tiles[x][y] = &Tile{false, false, false, false, false,"white", 0x1011, x, y}
				} else {
					m.Tiles[x][y] = &Tile{false, false, false, false, false,"white", 0x1019, x, y}
				}
			}
		}
	}
}

//черная магия
func (m *Map) GenerateRooms(src [][]int) {
	var predel1, predel2, predel3, predel4, ystart, yend int
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			src[y][x] = 12
		}
	}
	razrez1:= rand.Intn(m.Width/3) + 5
	razrez2:= rand.Intn(m.Width/3) + 3 
	razrez3:= m.Width - rand.Intn(m.Width/4) - razrez1
	razrez4:= m.Width - rand.Intn(m.Width/4) - 5
	razrez5:= rand.Intn(m.Height/2) + 3
	razrez6:= rand.Intn(m.Height/2) + (m.Height / 3)
	razrez7:= rand.Intn(m.Height/2) + (m.Height / 4)
	razrez8:= rand.Intn(m.Height/2) + (m.Height / 4) + 3
	for x := 0; x < razrez1; x++ {
		for y := 0; y < razrez5; y++ {
		  if x == 0 || x == razrez1 - 1 || y == 0 || y == razrez5 - 1 {
				src[y][x] = 11
			} else {
				src[y][x] = 0
			}
		}
	}

	for x := 0; x < razrez2; x++ {
		for y := razrez6; y < m.Height; y++ {
			if x == 0 || x == razrez2 - 1 || y == razrez6 || y == m.Height - 1 {
				src[y][x] = 11
			} else {
				src[y][x] = 0
			}
		}
	}

	for x := razrez3; x < m.Width; x++ {
		for y := 0; y < razrez7; y++ {
		  if x == razrez3 || x == m.Width - 1 || y == 0 || y == razrez7 - 1 {
				src[y][x] = 11
			} else {
				src[y][x] = 0
			}
		}
	}

	for x := razrez4; x < m.Width; x++ {
		for y := razrez8; y < m.Height; y++ {
			if x == razrez4 || x == m.Width - 1 || y == razrez8 || y == m.Height - 1 {
				src[y][x] = 11
			} else {
				src[y][x] = 0
			}
		}
	}

	if razrez1 < razrez2 {
		predel1 = razrez1 - 2 -  rand.Intn(razrez1/3)
	} else {
		predel1 = razrez2 - 2 -  rand.Intn(razrez2/3)
	}
	
	if predel1 < 1{
		predel1 = 1
	}
	
	if razrez3 < razrez4 {
		predel2 = razrez4 + 2 + rand.Intn(razrez3/5)
	} else {
		predel2 = razrez3 + 2 + rand.Intn(razrez4/5)
	}
	
	if predel2 > (m.Width - 2){
		predel2 = m.Width - 2
	}
	
	if razrez5 < razrez7 {
		predel3 = razrez5 - 2 - rand.Intn(razrez5/3)
	} else {
		predel3 = razrez7 - 2 - rand.Intn(razrez7/3)
	}
	
	if predel3 < 1{
		predel3 = 1
	}
	
	if razrez6 < razrez8 {
		predel4 = razrez8 + 1 + rand.Intn(razrez6/5)
	} else {
		predel4 = razrez6 + 1 + rand.Intn(razrez8/5)
	}
	
	if predel4 > (m.Height - 2){
		predel4 = m.Height - 2
	}


	if razrez5 < razrez6 +1 {
		ystart = razrez5 - 1
		yend = razrez6 + 1
	} else {
		ystart = razrez6
		yend = razrez6 + 1
	}
	for y := ystart; y < yend; y++ {
		src[y][predel1] = 0
		src[y][predel1 - 1] = 11
		src[y][predel1 + 1] = 11
	}
	
	if razrez7 < razrez8 + 1 {
		ystart = razrez7 - 1
		yend = razrez8 + 1
	} else {
		ystart = razrez8
		yend = razrez8 + 1
	}
	for y := ystart; y < yend; y++ {
		src[y][predel2] = 0
		src[y][predel2 - 1] = 11
		src[y][predel2 + 1] = 11
	}
	
	if razrez1 < razrez3 + 1 {
		ystart = razrez1 - 1
		yend = razrez3 + 1
	} else {
		ystart = razrez3
		yend = razrez3 + 1
	}
	for y := ystart; y < yend; y++ {
		src[predel3][y] = 0
		src[predel3 - 1][y] = 11
		src[predel3 + 1][y] = 11
	}
	
	if razrez2 < razrez4 + 1 {
		ystart = razrez2 - 1
		yend = razrez4 + 1
	} else {
		ystart = razrez4
		yend = razrez4 + 1
	}
	for y := ystart; y < yend; y++ {
		src[predel4][y] = 0
		src[predel4 - 1][y] = 11
		src[predel4 + 1][y] = 11
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
