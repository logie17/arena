package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Id      int
	Message chan string
	X       int
	Y       int
	EX      int
	EY      int
	HPLevel int
}

type Server struct {
	Clients        []*Client
	serverListener net.Listener
}

func NewClient(Id int) *Client {
	ch := make(chan string)
	rand.Seed(time.Now().Unix())
	x := rand.Intn(32-3) + 3
	y := rand.Intn(32-3) + 3

	return &Client{Id: Id, Message: ch, X: x, Y: y, HPLevel: 5}
}

func (client *Client) Listen(conn net.Conn) {
	go func() {
		for line := range client.Message {
			action, id, x, y := client.parseLine(string(line))
			if action == "stab" && client.NearEnemy() {
				client.sendAttackMsg(conn, action, id, x, y)
			} else {
				client.sendPosMsg(conn, action, id, x, y)
			}
		}
	}()
}

func (client *Client) sendPosMsg(conn net.Conn, action string, id, x, y int) {
	if id == client.Id {
		client.X = x
		client.Y = y
	} else {
		client.EX = x
		client.EY = y
	}
	client.SendMessage(conn, fmt.Sprintf("%s,%d,%d,%d\n", action, id, x, y))
}

func (client *Client) sendAttackMsg(conn net.Conn, action string, id, x, y int) {
	client.HPLevel--
	if client.HPLevel == 0 {
		client.SendMessage(conn, fmt.Sprintf("die,%d,%d,%d\n", id, x, y))
	} else {
		client.SendMessage(conn, fmt.Sprintf("hit,%d,%d,%d\n", id, x, y))
	}
}

func (client *Client) parseLine(line string) (string, int, int, int) {
	str := strings.Split(strings.TrimSpace(line), ",")
	action := str[0]
	iid, _ := strconv.Atoi(str[1])
	x, _ := strconv.Atoi(str[2])
	y, _ := strconv.Atoi(str[3])
	return action, iid, x, y
}

func (client *Client) NearEnemy() bool {
	if math.Abs(float64(client.X-client.EX)) <= 1 && math.Abs(float64((client.Y-client.EY))) <= 1 {
		return true
	} else {
		return false
	}
}

func (client *Client) SendMessage(conn net.Conn, msg string) {
	conn.Write([]byte(msg))
}

func NewServer() *Server {
	s := new(Server)
	return s

}

func (s *Server) Serve() (err error) {
	s.serverListener, err = net.Listen("tcp", fmt.Sprint("127.0.0.1:", 9000))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer s.serverListener.Close()
	connId := 1
	for {
		//This blocks
		conn, err := s.serverListener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		client := NewClient(connId)
		client.Listen(conn)
		s.Clients = append(s.Clients, client)
		go s.handleConn(conn, client)
		connId++
	}
	return nil

}

func (s *Server) handleConn(conn net.Conn, client *Client) {
	fmt.Println("trying to handle connection")
	bufc := bufio.NewReader(conn)
	s.Broadcast(fmt.Sprintf("pos,%d,%d,%d\n", client.Id, client.X, client.Y))
	for {
		line, _, err := bufc.ReadLine()
		if err != nil {
			break
		}
		s.Broadcast(string(line))
	}
}

func (s *Server) Broadcast(line string) {
	fmt.Println(line)
	for _, client := range s.Clients {
		client.Message <- string(line)
	}
}

func main() {
	s := NewServer()
	err := s.Serve()
	if err != nil {
		fmt.Println(err)
	}
}
