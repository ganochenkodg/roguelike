package entity

import (
	blt "bearlibterminal"
	"dijkstramaps"
	"gamemap"
	"math/rand"
	"strconv"
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
}

func (e *GameEntity) Move(dx int, dy int) {
	// Move the entity by the amount (dx, dy)
	e.X += dx
	e.Y += dy
}

func (e *GameEntity) Draw() {
	// Draw the entity to the screen
	var hpbar int
	blt.Layer(e.Layer)
	blt.Color(blt.ColorFromName(e.Color))
  blt.Put(e.X*4, e.Y*2, e.Char)
	blt.Layer(20)
	blt.Color(blt.ColorFromName("white"))
	
	if hpbar = 7 - (e.HP[0] * 6) / e.HP[1]; hpbar > 6{
		hpbar = 6
	}
	blt.Put(e.X*4, e.Y*2, hpbar + 0x3000)
	
}

func (e *GameEntity) Clear() {
	// Remove the entity from the screen
	blt.Layer(e.Layer)
	blt.Put(e.X*4, e.Y*2, 0)
	blt.Layer(20)
	blt.Put(e.X*4, e.Y*2, 0)
}

func (e *GameEntity) Hunting(edm *dijkstramaps.EntityDijkstraMap, gamemap *gamemap.Map, player *GameEntity) {
	var oldx, oldy, newx, newy = e.X, e.Y, 0, 0
	// Check to see if the provided coordinates contain a blocked tile
	if e.NPC && edm.ValuesMap[e.X][e.Y] == 1 {
		blt.Print(1,43,e.Fight(gamemap, player))
	}
	if e.NPC && edm.ValuesMap[e.X][e.Y] < 9 && edm.ValuesMap[e.X][e.Y] > 1{
		for x := -1; x < 2; x++ {
	 		for y := -1; y < 2; y++ {
				if edm.ValuesMap[oldx+x][oldy+y] < edm.ValuesMap[oldx][oldy] && !gamemap.Tiles[oldx+x][oldy+y].Mob {
          newx, newy = x, y
				}
			}
		}
	e.Move(newx,newy)
//	blt.Refresh()
	}
}

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
			target.Clear()
			gamemap.Tiles[target.X][target.Y].Mob = false
			result = result + ".You kill " + target.Name
		}
	}
	return result
}
