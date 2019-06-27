package messages

import (
  blt "bearlibterminal"
)

type Messages []string

func (m Messages) Initialize() {
	for x:=0; x < 5; x++ {
		m[x] = ""
	}
}

func (m Messages) AddMessage(s string) {
	for x:=0; x < 4; x++ {
		m[x] = m[x+1]
	}
	m[4] = s
}

func (m Messages) PrintMessages() {
	blt.Layer(5)
	blt.ClearArea(1,42,100,5)
	for x:=0; x < 5; x++ {
		blt.Print(1,41+x,m[x])
	}
}
