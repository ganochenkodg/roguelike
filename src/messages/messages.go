package messages

import (
  blt "bearlibterminal"
)

type Messages []string

func (m Messages) Initialize() {
	for x:=0; x < 30; x++ {
		m[x] = ""
	}
}

func (m Messages) AddMessage(s string) {
	for x:=0; x < 29; x++ {
		m[x] = m[x+1]
	}
	m[29] = s
}

func (m Messages) PrintMessages() {
	blt.Layer(5)
	blt.ClearArea(1,32,100,5)
	for x:=0; x < 5; x++ {
		blt.Print(1,31+x,m[x+24])
	}
}

func (m Messages) DrawJournal() {
  blt.Clear()
	blt.Layer(5)
  blt.Print(1,1,"Game journal:")
	for x:=0; x < 30; x++ {
		blt.Print(1,3+x,m[x])
	}
  blt.Print(1,34,"Press J or ESC to return")
}
