package main

import (
	"bufio"
	"fmt"
	"github.com/logie17/arena/client/fighter"
	"github.com/nsf/termbox-go"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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

func readFromServer(fighterId int, fighters []fighter.Fighter, bufc *bufio.Reader, reply chan fighter.CommandData) {
	go func() {
		for {
			line, _ := bufc.ReadString('\n')
			str := strings.Split(strings.TrimSpace(string(line)), ",")
			id, _ := strconv.Atoi(str[1])

			if id != fighterId && isNewEnemy(id, fighters) {
				x, _ := strconv.Atoi(str[2])
				y, _ := strconv.Atoi(str[3])
				enemy := fighter.NewFighter(x, y, id, "enemy", reply)
				fighters = append(fighters, enemy)
			}

			for _, fighter := range fighters {
				fighter.SendMessage(line)
			}
		}
	}()
}

func isNewEnemy(id int, fighters []fighter.Fighter) bool {
	isNew := true
	for _, fighter := range fighters {
		if id == fighter.Id() {
			isNew = false
		}
	}
	return isNew
}

func handleFighterActions(cn net.Conn, reply chan fighter.CommandData) {
	go func() {
		for {
			select {
			case response := <-reply:
				val := response.Value
				action := response.Action
				id := val[0]
				x := val[1]
				y := val[2]

				if action == "FLUSH" {
					termbox.Flush()
				} else if action == "HIT" {
					termbox.SetCell(x, y, '@', termbox.ColorYellow, termbox.ColorBlack)
					termbox.Flush()
					go func() {
						time.Sleep(100 * time.Millisecond)
						termbox.SetCell(x, y, '@', termbox.ColorRed, termbox.ColorBlack)
						termbox.Flush()

					}()
				} else if action == "HIDE" {
					termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
					termbox.Flush()
				} else if action == "DRAW" {
					enemy := val[3]
					if enemy == 1 {
						termbox.SetCell(x, y, '@', termbox.ColorRed, termbox.ColorBlack)
					} else {
						termbox.SetCell(x, y, '@', termbox.ColorBlue, termbox.ColorBlack)
					}
					termbox.Flush()
				} else if action == "KILL" {
					println("YOU DIED!")
					os.Exit(1)
				} else if action == "WIN" {
					println("YOU WIN!!!")
					os.Exit(1)
				} else {
					cn.Write([]byte(fmt.Sprintf("%s,%d,%d,%d\n", action, id, x, y)))
				}

			}
		}
	}()

}

func setupBoard() {
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printMsg(int(boardWidth/2)-(int(boardWidth/2)/2), 0, termbox.ColorRed, termbox.ColorBlack, "ARENA!!! FIGHT!!!")

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

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	defer termbox.Close()
	setupBoard()

	destination := "127.0.0.1:9000"
	cn, err := net.Dial("tcp", destination)
	if err != nil {
		os.Exit(1)
	}

	//defer cn.Close();

	bufc := bufio.NewReader(cn)
	if err != nil {
		os.Exit(1)
	}

	line, err := bufc.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}

	str := strings.Split(strings.TrimSpace(string(line)), ",")

	fighterId, err := strconv.Atoi(str[1])
	if err != nil {
		os.Exit(1)
	}
	x, err := strconv.Atoi(str[2])
	if err != nil {
		os.Exit(1)
	}
	y, err := strconv.Atoi(str[3])
	if err != nil {
		os.Exit(1)
	}

	reply := make(chan fighter.CommandData, 4)
	player := fighter.NewFighter(x, y, fighterId, "me", reply)
	fighters := []fighter.Fighter{player}

	readFromServer(fighterId, fighters, bufc, reply)
	handleFighterActions(cn, reply)
	handleKeyEvents(player)
}

func handleKeyEvents(f fighter.Fighter) {
LOOP:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break LOOP
			case termbox.KeyArrowDown:
				f.Action("Down")
			case termbox.KeyArrowUp:
				f.Action("Up")
			case termbox.KeyArrowLeft:
				f.Action("Left")
			case termbox.KeyArrowRight:
				f.Action("Right")
			case termbox.KeySpace:
				f.Action("Stab")
			}
		}
	}
}
