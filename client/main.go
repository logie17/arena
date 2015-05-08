package main

import (
	"github.com/nsf/termbox-go"
	"net"
	"strings"
	"strconv"
	"fmt"
	"bufio"
	"time"
	"github.com/logie17/arena/safehash"
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

var mySafeMap = safehash.NewSafeMap()

type fighter struct {
	x int
	y int
	id int
	enemyx int
	enemyy int
	enemyid int
	kind string
	name string
	character rune
	message chan string
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
	Stab()
	Hit()
	SetEnemyId(int)
	Listen()
	SendMessage(string)
}

func (fighter *fighter) SendMessage(line string) {
	fighter.message<-line
}

func (fighter * fighter) Id() int {
	return fighter.id
}

func (fighter * fighter) SetId(id int){
	fighter.id = id
}

func (fighter * fighter) SetEnemyId(id int){
	fighter.enemyid = id
}

func NewFighter(x, y, id int, kind string, conn net.Conn) Fighter {
	message := make(chan string)
	fighter := &fighter{x, y, id, 0, 0, 0, kind, "Bad ass", '@', message, conn}
	fighter.Listen()
	fighter.Draw()
	termbox.Flush()

	return fighter
}

func (fighter * fighter) Pos(x, y int) {
	fighter.Hide()
	fighter.x = x
	fighter.y = y
	fighter.Draw()
}

func (fighter * fighter) Action(action string) {
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
		fighter.Stab()
	}

	fighter.conn.Write([]byte(fmt.Sprintf("%s,%d,%d,%d\n",act,fighter.id,fighter.x,fighter.y)))
}

func (fighter *fighter) Listen() {
	go func() {
		for line := range fighter.message {
			str := strings.Split(strings.TrimSpace(string(line)),",");
			action := str[0]
			id,_ := strconv.Atoi(str[1])

			if id == fighter.id && fighter.kind == "enemy" && action == "pos" {
				x,_ := strconv.Atoi(str[2])
				y,_ := strconv.Atoi(str[3])
				mySafeMap.Insert(fmt.Sprintf("%d_x",id),x)
				mySafeMap.Insert(fmt.Sprintf("%d_y",id),y)
				fighter.Pos(x,y)
				termbox.Flush()
			}
			
			if  id != fighter.id {
 				fighter.enemyid = id
			}

			if action == "hit" && id != fighter.id {
				fighter.Hit()
			}
		}
	}()

}

func (fighter * fighter) Stab() {
//	fmt.Println("STAB!!!")
}

func (fighter * fighter) Hit() {
	termbox.SetCell(fighter.x, fighter.y, fighter.character, termbox.ColorYellow, termbox.ColorBlack)
	termbox.Flush()
	go func() {
		time.Sleep(100 * time.Millisecond)
		termbox.SetCell(fighter.x, fighter.y, fighter.character, termbox.ColorRed, termbox.ColorBlack)
		termbox.Flush()

	}()
}

func (fighter * fighter) Down() {
	fighter.Hide()
	newY := fighter.y + 1
	if fighter.y < 33  && !fighter.cellIsOccupied(fighter.x, newY) {
		fighter.y = newY
	}
	fighter.Draw()
}

func (fighter * fighter) Up() {
	fighter.Hide()
	newY := fighter.y - 1
	if fighter.y > 3  && !fighter.cellIsOccupied(fighter.x, newY) {
		fighter.y = newY
	}
	fighter.Draw()
}

func (fighter *fighter) cellIsOccupied(x,y int) bool {
	enemyPosX := mySafeMap.Find(fmt.Sprintf("%d_x",fighter.enemyid))
	enemyPosY := mySafeMap.Find(fmt.Sprintf("%d_y",fighter.enemyid))
	if y == enemyPosY && x == enemyPosX {
		return true
	}
	return false
}

func (fighter * fighter) Right() {
	fighter.Hide()
	newX := fighter.x + 1
	if fighter.x < 80  && !fighter.cellIsOccupied(newX, fighter.y) {
		fighter.x = newX
	}
	fighter.Draw()
}

func (fighter * fighter) Left() {
	fighter.Hide()
	newX := fighter.x - 1
	if fighter.x > 0 && !fighter.cellIsOccupied(newX, fighter.y) {
		fighter.x = newX
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

	fighterId,_ := strconv.Atoi(str[1])
	x,_ := strconv.Atoi(str[2])
	y,_ := strconv.Atoi(str[3])

	mySafeMap.Insert(fmt.Sprintf("%d_x",fighterId),x)
	mySafeMap.Insert(fmt.Sprintf("%d_y",fighterId),y)

	fighter := NewFighter(x,y,fighterId,"me",cn)
	fighters := []Fighter{fighter}
	var enemy Fighter

	go func() {
		for {
			line, _ := bufc.ReadString('\n')
			str := strings.Split(strings.TrimSpace(string(line)),",");
			id,_ := strconv.Atoi(str[1])

			if id != fighterId && enemy == nil {
				x,_ := strconv.Atoi(str[2])
				y,_ := strconv.Atoi(str[3])
				enemy = NewFighter(x,y,id,"enemy",cn)
				fighters = append(fighters, enemy)
			}

			for _, fighter := range fighters {
				fighter.SendMessage(line)
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
			case termbox.KeySpace:
				fighter.Action("Stab")
				termbox.Flush()

			}
		}
	}

}
