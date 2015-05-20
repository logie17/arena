package fighter

import (
	"fmt"
	"github.com/logie17/arena/safehash"
)

type Line struct {
	Action string
	Id     int
	X      int
	Y      int
}

type fighter struct {
	x         int
	y         int
	id        int
	enemyx    int
	enemyy    int
	enemyid   int
	kind      string
	character rune
	message   chan Line
	reply     chan Line
}

type Fighter interface {
	Left()
	Right()
	Up()
	Down()
	Id() int
	Action(string)
	Listen()
	SendMessage(Line)
}

var mySafeMap = safehash.NewSafeMap()

func (fighter *fighter) SendMessage(line Line) {
	fighter.message <- line
}

func (fighter *fighter) Id() int {
	return fighter.id
}

func NewFighter(x, y, id int, kind string, reply chan Line) Fighter {
	mySafeMap.Insert(fmt.Sprintf("%d_x", id), x)
	mySafeMap.Insert(fmt.Sprintf("%d_y", id), y)
	message := make(chan Line)
	fighter := &fighter{
		x: x, y: y, id: id, kind: kind, character: '♥', //code point 2665
		message: message, reply: reply,
	}
	fighter.Listen()
	fighter.Draw()
	return fighter
}

func (fighter *fighter) Pos(x, y int) {
	fighter.Hide()
	fighter.x = x
	fighter.y = y
	fighter.Draw()
}

func (fighter *fighter) Action(action string) {
	act := "pos"
	switch action {
	case "Down":
		fighter.Down()
	case "Up":
		fighter.Up()
	case "Left":
		fighter.Left()
	case "Right":
		fighter.Right()
	case "Stab":
		act = "stab"
	}

	fighter.reply <- Line{act, fighter.id, fighter.x, fighter.y}
}

func (fighter *fighter) Listen() {
	go func() {
		for line := range fighter.message {
			action := line.Action
			id := line.Id
			x := line.X
			y := line.Y

			if id == fighter.id && fighter.kind == "enemy" && action == "pos" {
				mySafeMap.Insert(fmt.Sprintf("%d_x", id), x)
				mySafeMap.Insert(fmt.Sprintf("%d_y", id), y)
				fighter.Pos(x, y)
				fighter.reply <- Line{"refresh_board", id, x, y}
			}

			if id != fighter.id {
				fighter.enemyid = id
			}

			if action == "hit" && id != fighter.id {
				fighter.reply <- Line{"hit", fighter.id, fighter.x, fighter.y}
			}

			if id != fighter.id && action == "die" {
				fighter.reply <- Line{"win", fighter.id, fighter.x, fighter.y}
			}

			if id == fighter.id && action == "die" {
				fighter.reply <- Line{"kill", fighter.id, fighter.x, fighter.y}
			}
		}
	}()
}

func (fighter *fighter) Down() {
	fighter.Hide()
	newY := fighter.y + 1
	if fighter.y < 33 && !fighter.cellIsOccupied(fighter.x, newY) {
		fighter.y = newY
	}
	fighter.Draw()
}

func (fighter *fighter) Up() {
	fighter.Hide()
	newY := fighter.y - 1
	if fighter.y > 3 && !fighter.cellIsOccupied(fighter.x, newY) {
		fighter.y = newY
	}
	fighter.Draw()
}

func (fighter *fighter) Right() {
	fighter.Hide()
	newX := fighter.x + 1
	if fighter.x < 80 && !fighter.cellIsOccupied(newX, fighter.y) {
		fighter.x = newX
	}
	fighter.Draw()
}

func (fighter *fighter) Left() {
	fighter.Hide()

	newX := fighter.x - 1
	if fighter.x > 0 && !fighter.cellIsOccupied(newX, fighter.y) {
		fighter.x = newX
	}
	fighter.Draw()
}

func (fighter *fighter) cellIsOccupied(x, y int) bool {
	enemyPosX := mySafeMap.Find(fmt.Sprintf("%d_x", fighter.enemyid))
	enemyPosY := mySafeMap.Find(fmt.Sprintf("%d_y", fighter.enemyid))
	if y == enemyPosY && x == enemyPosX {
		return true
	}
	return false
}

func (fighter *fighter) Hide() {
	fighter.reply <- Line{"hide", fighter.id, fighter.x, fighter.y}
}

func (fighter *fighter) Draw() {
	if fighter.kind == "enemy" {
		fighter.reply <- Line{"redraw_enemy", fighter.id, fighter.x, fighter.y}
	} else {
		fighter.reply <- Line{"redraw_me", fighter.id, fighter.x, fighter.y}
	}
}
