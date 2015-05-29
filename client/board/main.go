package board

import (
	"github.com/nsf/termbox-go"
	"net"
	"os"
)

const (
	boardWidth  = 79
	boardHeight = 30
)

func printMsg(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func UpdateCell(x, y int, c rune, color termbox.Attribute) {
	termbox.SetCell(x, y, c, color, termbox.ColorBlack)
	termbox.Flush()
}

func DrawBoard(msg string) {
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printMsg(int(boardWidth/2)-(int(boardWidth/2)/2), 0, termbox.ColorRed, termbox.ColorBlack, msg)

	for i := 1; i < 80; i++ {
		termbox.SetCell(i, 2, 0x2500, termbox.ColorGreen, termbox.ColorBlack)
		termbox.SetCell(i, 33, 0x2500, termbox.ColorGreen, termbox.ColorBlack)
	}

	for i := 2; i < 33; i++ {
		termbox.SetCell(0, i, 0x2502, termbox.ColorGreen, termbox.ColorBlack)
		termbox.SetCell(80, i, 0x2502, termbox.ColorGreen, termbox.ColorBlack)
	}

	termbox.Flush()

}

func InitBoard() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	DrawBoard("ARENA!! FIGHT TO THE DEATH!!")

}

func Close() {
	termbox.Close()
}

type Fighter interface {
	Action(termbox.Key) []byte
}

func HandleKeyEvents(cn net.Conn, f Fighter) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				os.Exit(0)
			default:
				cn.Write(f.Action(ev.Key))
			}
		}
	}
}
