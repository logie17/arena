package fighter

import (
	"fmt"
	"github.com/logie17/arena/client/board"
	"github.com/logie17/arena/safehash"
	"github.com/nsf/termbox-go"
	"time"
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
	X() int
	Y() int
}

var mySafeMap = safehash.NewSafeMap()

func (fighter *fighter) SendMessage(line Line) {
	fighter.message <- line
}

func (fighter *fighter) Id() int {
	return fighter.id
}

func (f *fighter) X() int {
	return f.x
}

func (f *fighter) Y() int {
	return f.y
}

func NewFighter(x, y, id int, kind string) Fighter {
	mySafeMap.Insert(fmt.Sprintf("%d_x", id), x)
	mySafeMap.Insert(fmt.Sprintf("%d_y", id), y)
	message := make(chan Line)
	fighter := &fighter{
		x: x, y: y, id: id, kind: kind, character: '♥', //code point 2665
		message: message,
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
	switch action {
	case "Down":
		fighter.Down()
	case "Up":
		fighter.Up()
	case "Left":
		fighter.Left()
	case "Right":
		fighter.Right()
	}

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
			}

			if id != fighter.id {
				fighter.enemyid = id
			}

			if action == "hit" && id != fighter.id {
				board.UpdateCell(x, y, '♥', termbox.ColorYellow)
				go func() {
					time.Sleep(100 * time.Millisecond)
					board.UpdateCell(x, y, '♥', termbox.ColorRed)

				}()

			}

			if id != fighter.id && action == "die" {
				board.DrawBoard("YOU WIN - GAME OVER!")
			}

			if id == fighter.id && action == "die" {
				board.DrawBoard("YOU DIED!! - GAME OVER!")
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
	termbox.SetCell(fighter.x, fighter.y, ' ', termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()

}

func (f *fighter) Draw() {
	if f.kind == "enemy" {
		board.UpdateCell(f.x, f.y, '♥', termbox.ColorRed)
	} else {
		board.UpdateCell(f.x, f.y, '♥', termbox.ColorCyan)
	}
}
