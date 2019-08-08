package entity

import (
	blt "bearlibterminal"
	"dijkstramaps"
	"gamemap"
	"math/rand"
	"strconv"
	"messages"
)

type GameEntity struct {
	X int
	Y int
	Layer int
	Char int
	Color string
	NPC bool
	Name string
	HP []int
	Vision int
	Speed float32
	SpeedPool float32
}

//переместить entity на новую позицию
func (e *GameEntity) Move(dx int, dy int) {
	e.X += dx
	e.Y += dy
}

//отрисовать entity
func (e *GameEntity) Draw(camerax, cameray int) {
	var hpbar int
	blt.Layer(e.Layer)
	blt.Color(blt.ColorFromName(e.Color))
	newx, newy := GetCamera(e.X, e.Y, camerax, cameray)
  blt.Put(newx*4, newy*2, e.Char)
//рисуем полоску hp. надо бы сделать больше шагов бара, хотяб 10
	blt.Layer(4)
	blt.Color(blt.ColorFromName("white"))
	
	if hpbar = 7 - (e.HP[0] * 6) / e.HP[1]; hpbar > 6{
		hpbar = 6
	}
	blt.Put(newx*4, newy*2, hpbar + 0x3000)
	
}

//рисуем на позиции пустой символ
func (e *GameEntity) Clear(camerax, cameray int) {
	newx, newy := GetCamera(e.X, e.Y, camerax, cameray)
	blt.Layer(e.Layer)
	blt.Put(newx*4, newy*2, 0)
	blt.Layer(4)
	blt.Put(newx*4, newy*2, 0)
}

//ежеходовая проверка не пора ли преследовать игрока
func (e *GameEntity) Hunting(edm *dijkstramaps.EntityDijkstraMap, gamemap *gamemap.Map, player *GameEntity, message messages.Messages) {
	var oldx, oldy, newx, newy = e.X, e.Y, 0, 0
	if e.NPC && edm.ValuesMap[e.X][e.Y] < e.Vision{
		e.SpeedPool += (100.0 / player.Speed)
	} else {
		e.SpeedPool = 0.0
	}
  //если проверяем нпс и стоит около игрока на начало хода то вызываем сражение
	for e.SpeedPool > (100.0 / e.Speed) {
		e.SpeedPool -= (100.0 / e.Speed)
	if e.NPC && edm.ValuesMap[e.X][e.Y] == 1 {
		message.AddMessage(e.Fight(gamemap, player))
	}
	//если нпс, дальше чем на одну клетку и игрок в радиусе видимости то пытаемся подойти
	if e.NPC && edm.ValuesMap[e.X][e.Y] < e.Vision && edm.ValuesMap[e.X][e.Y] > 1{
		for x := -1; x < 2; x++ {
	 		for y := -1; y < 2; y++ {
				//ищем клетку с меньшим весом и без мобов
				if edm.ValuesMap[oldx+x][oldy+y] < edm.ValuesMap[oldx][oldy] && !gamemap.Tiles[oldx+x][oldy+y].Mob {
          newx, newy = x, y
				}
			}
		}
	e.Move(newx,newy)
	oldx, oldy, newx, newy = e.X, e.Y, 0, 0
	}
}
}

//битва двух entity, пока что тупорылая до безумия
func (e *GameEntity) Fight(gamemap *gamemap.Map, target *GameEntity) string{
	var result string
  kick:=rand.Intn(5)
	target.HP[0]-=kick
	blt.Layer(1)
	blt.Color(blt.ColorFromName("white"))
	result = e.Name + " kicks " + target.Name + " for " + strconv.Itoa(kick)
	if target.HP[0] < 1{
		if e.NPC {
			result = result + ". You lose!"
		} else {
			target.Clear(gamemap.CameraX, gamemap.CameraY)
			gamemap.Tiles[target.X][target.Y].Mob = false
			result = result + ".You kill " + target.Name
		}
	}
	return result
}

func GetCamera(x ,y, camerax, cameray int) (int, int){
	newx := 12 + x - camerax
	newy := 8 + y - cameray
	if newy > 14{
		newy = -1
		} 
	return newx, newy
}
