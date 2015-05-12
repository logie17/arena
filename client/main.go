package main

import (
	"github.com/nsf/termbox-go"
	"net"
	"strings"
	"strconv"
	"fmt"
	"bufio"
	"time"
	"github.com/logie17/arena/client/fighter"

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
	termbox.Flush()

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

	reply := make(chan fighter.CommandData,4)
	fighter1 := fighter.NewFighter(x,y,fighterId,"me",cn, reply)

	fighters := []fighter.Fighter{fighter1}
	var enemy fighter.Fighter

	go func() {
		for {
			line, _ := bufc.ReadString('\n')
			str := strings.Split(strings.TrimSpace(string(line)),",");
			id,_ := strconv.Atoi(str[1])

			if id != fighterId && enemy == nil {
				x,_ := strconv.Atoi(str[2])
				y,_ := strconv.Atoi(str[3])
				enemy = fighter.NewFighter(x,y,id,"enemy",cn, reply)
				fighters = append(fighters, enemy)
			}

			for _, fighter := range fighters {
				fighter.SendMessage(line)
			}
		}
	}()

	go func() {
		for {
			select {
			case response := <-reply:
				val := response.Value
				action := response.Action
				x := val[0]
				y := val[1]


				if (action == "FLUSH") {
					termbox.Flush()
				} else if ( action == "HIT") {
					termbox.SetCell(x, y, '@', termbox.ColorYellow, termbox.ColorBlack)
					termbox.Flush()
					go func() {
						time.Sleep(100 * time.Millisecond)
						termbox.SetCell(x, y, '@', termbox.ColorRed, termbox.ColorBlack)
						termbox.Flush()

					}()
				} else if ( action == "HIDE" ) {
					termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
					termbox.Flush()
				} else if ( action == "DRAW" ) {
					enemy := val[2]
					if enemy == 1 {
						termbox.SetCell(x, y, '@', termbox.ColorRed, termbox.ColorBlack)
					} else {
						termbox.SetCell(x, y, '@', termbox.ColorBlue, termbox.ColorBlack)
					}
					termbox.Flush()
				}
				
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
				fighter1.Action("Down")
				termbox.Flush()
			case termbox.KeyArrowUp:
				fighter1.Action("Up")
				termbox.Flush()
			case termbox.KeyArrowLeft:
				fighter1.Action("Left")
				termbox.Flush()
			case termbox.KeyArrowRight:
				fighter1.Action("Right")
				termbox.Flush()
			case termbox.KeySpace:
				fighter1.Action("Stab")
				termbox.Flush()

			}
		}
	}

}
