package main

import (
	"bufio"
	"fmt"
	"github.com/logie17/arena/client/board"
	"github.com/logie17/arena/client/fighter"
	"net"
	"os"
	"strconv"
	"strings"
)

func readFromServer(fighterId int, fighters []fighter.Fighter, bufc *bufio.Reader) {
	go func() {
		for {
			line, _ := bufc.ReadString('\n') // Blocks
			data := parseLine(string(line))

			if data.Id != fighterId && isNewEnemy(data.Id, fighters) {
				enemy := fighter.NewFighter(data.X, data.Y, data.Id, "enemy")
				fighters = append(fighters, enemy)
			}

			for _, fighter := range fighters {
				fighter.SendMessage(data)
			}
		}
	}()
}

func parseLine(line string) fighter.Line {
	str := strings.Split(strings.TrimSpace(string(line)), ",")
	action := str[0]
	id, _ := strconv.Atoi(str[1])
	x, _ := strconv.Atoi(str[2])
	y, _ := strconv.Atoi(str[3])

	return fighter.Line{action, id, x, y}
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

func establishConnection() net.Conn {
	destination := "127.0.0.1:9000"

	cn, err := net.Dial("tcp", destination)
	if err != nil {
		fmt.Println("Unable to open connection: ", err.Error())
		os.Exit(1)
	}
	return cn
}

func readConnectionLine(bufc *bufio.Reader) (int, int, int) {
	line, err := bufc.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to read connection string", err.Error())
		os.Exit(1)
	}

	data := parseLine(string(line))
	return data.X, data.Y, data.Id
}

func main() {
	board.InitBoard()
	defer board.Close()
	cn := establishConnection()
	defer cn.Close()

	bufc := bufio.NewReader(cn)

	x, y, fighterId := readConnectionLine(bufc)
	player := fighter.NewFighter(x, y, fighterId, "me")
	fighters := []fighter.Fighter{player}

	readFromServer(fighterId, fighters, bufc)
	board.HandleKeyEvents(cn, player)
}
