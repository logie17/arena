package main

import (
	"github.com/nsf/termbox-go"
	"net"
	"strings"
	"strconv"
	"fmt"
	"bufio"
)

const (
	boardWidth  = 79
	boardHeight = 30
)

func print_msg(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

type fighter struct {
	x int
	y int
	id int
	kind string
	name string
	character rune
	conn net.Conn
}

type Fighter interface {
	Draw()
	Hide()
	Left()
	Right()
	Up()
	Down()
	Id() int
	Pos(int, int)
	Action(string)
	SetId(int)
	// Stab()
}

func (fighter * fighter) Id() int {
	return fighter.id
}

func (fighter * fighter) SetId(id int){
	fighter.id = id
}

func NewFighter(x, y, id int, kind string, conn net.Conn) Fighter {
	return &fighter{x, y, id, kind, "Bad ass", '@', conn}
}

func (fighter * fighter) Pos(x, y int) {
	fighter.Hide()
	fighter.x = x
	fighter.y = y
}

func (fighter * fighter) Action(action string) {
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

	fighter.conn.Write([]byte(fmt.Sprintf("%d,%d,%d\n",fighter.id,fighter.x,fighter.y)))
}

func (fighter * fighter) Down() {
	fighter.Hide()
	if fighter.y < 33 {
		fighter.y++
	}
	fighter.Draw()
}

func (fighter * fighter) Up() {
	fighter.Hide()
	if fighter.y > 3 {
		fighter.y--
	}
	fighter.Draw()
}

func (fighter * fighter) Right() {
	fighter.Hide()
	if fighter.x < 80 {
		fighter.x++
	}
	fighter.Draw()
}

func (fighter * fighter) Left() {
	fighter.Hide()
	if fighter.x > 0 {
		fighter.x--
	}
	fighter.Draw()
}

func (fighter * fighter) Hide() {
	termbox.SetCell(fighter.x, fighter.y, ' ', termbox.ColorBlack, termbox.ColorBlack)
}

func (fighter * fighter) Draw() {
	if fighter.kind == "enemy" {
		termbox.SetCell(fighter.x, fighter.y, fighter.character, termbox.ColorRed, termbox.ColorBlack)
	} else {
		termbox.SetCell(fighter.x, fighter.y, fighter.character, termbox.ColorBlue, termbox.ColorBlack)
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	print_msg(int(boardWidth/2) - (int(boardWidth/2)/2), 0, termbox.ColorRed, termbox.ColorBlack, "ARENA!!! FIGHT!!!")
	
	// termbox.SetCell(0, 2, 0x250C, termbox.ColorRed, termbox.ColorBlack)
	// termbox.SetCell(boardWidth+1, 0, 0x2510, termbox.ColorRed, termbox.ColorBlack)
	// termbox.SetCell(0, boardHeight+1, 0x2514, termbox.ColorRed, termbox.ColorBlack)
	// termbox.SetCell(boardWidth+1, boardHeight+1, 0x2515, termbox.ColorRed, termbox.ColorBlack)

	for i := 1; i < 80; i++ {
		termbox.SetCell(i, 2, 0x2500, termbox.ColorRed, termbox.ColorBlack)
		termbox.SetCell(i, 31, 0x2500, termbox.ColorRed, termbox.ColorBlack)
	}

	for i := 2; i < 33; i++ {
		termbox.SetCell(0, i, 0x2502, termbox.ColorRed, termbox.ColorBlack)
		termbox.SetCell(80, i, 0x2502, termbox.ColorRed, termbox.ColorBlack)
	}

	// TODO Chat Box
	for i := 83; i < 120; i++ {
		termbox.SetCell(i, 2, 0x2500, termbox.ColorRed, termbox.ColorBlack)
		termbox.SetCell(i, 31, 0x2500, termbox.ColorRed, termbox.ColorBlack)
	}


	destination := "127.0.0.1:9000";
	cn, err := net.Dial("tcp", destination);
	//defer cn.Close();

	bufc := bufio.NewReader(cn)

	if err != nil {
		fmt.Println("YUCK!")
	}
	line, _ := bufc.ReadString('\n')
	str := strings.Split(strings.TrimSpace(string(line)),",");

	
	fmt.Println(str)

	id,_ := strconv.Atoi(str[0])
	x,_ := strconv.Atoi(str[1])
	y,_ := strconv.Atoi(str[2])

	fighter := NewFighter(x,y,id,"me",cn)
	fighter.Draw()
	termbox.Flush()

	enemy := NewFighter(0,0, 0,"enemy",cn)

	go func() {
		for {
			line, _ := bufc.ReadString('\n')
			str := strings.Split(strings.TrimSpace(string(line)),",");

			id,_ := strconv.Atoi(str[0])
			if id != fighter.Id() && enemy.Id() == 0 {
				enemy.SetId(id)
			}

			if id == enemy.Id() {
				x,_ := strconv.Atoi(str[1])
				y,_ := strconv.Atoi(str[2])
				enemy.Pos(x,y)
				enemy.Draw()
				termbox.Flush()
			}
		}
	}()
LOOP:
	for {
	
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break LOOP
			case termbox.KeyArrowDown:
				fighter.Action("Down")
				termbox.Flush()
			case termbox.KeyArrowUp:
				fighter.Action("Up")
				termbox.Flush()
			case termbox.KeyArrowLeft:
				fighter.Action("Left")
				termbox.Flush()
			case termbox.KeyArrowRight:
				fighter.Action("Right")
				termbox.Flush()

			}
		}
	}

}
