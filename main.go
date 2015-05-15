package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"math/rand"
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
	Fighter1Port   uint16
	Hostname       string
	Logger         *log.Logger
	Clients        []Client
	clientListener net.Listener
	serverListener net.Listener
}

func NewClient(Id int) *Client {
	ch := make(chan string)
	return &Client{Id, ch, 20, 20, 0, 0, 5}
}

func (client *Client) Listen(conn net.Conn) {
	go func() {
		for line := range client.Message {
			str := strings.Split(strings.TrimSpace(string(line)), ",")
			action := str[0]
			iid, _ := strconv.Atoi(str[1])
			x, _ := strconv.Atoi(str[2])
			y, _ := strconv.Atoi(str[3])
			fmt.Println(str)
			if action == "stab" && client.NearEnemy() {
				client.HPLevel--
				if (client.HPLevel == 0 ) {
					client.SendMessage(conn, fmt.Sprintf("die,%d\n", iid))
				} else {
					client.SendMessage(conn, fmt.Sprintf("hit,%d\n", iid))
				}

			} else {
				if iid == client.Id {
					client.X = x
					client.Y = y
				} else {
					client.EX = x
					client.EY = y
				}
				client.SendMessage(conn, fmt.Sprintf("%s,%d,%d,%d\n", action, iid, x, y))
			}
		}
	}()
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
	s.Fighter1Port = 9000
	s.Hostname = "localhost"

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		names, _ := net.LookupAddr(addr.String())
		if len(names) > 0 {
			s.Hostname = names[0]
			break
		}
	}

	s.Clients = []Client{}
	return s

}

func (s *Server) log(v interface{}) {
	if s.Logger != nil && v != nil {
		s.Logger.Println(v)
	}
}

func (s *Server) Serve() (err error) {
	s.serverListener, err = net.Listen("tcp", fmt.Sprint("127.0.0.1:", s.Fighter1Port))
	if err != nil {
		s.log(err)
		return err
	}

	defer s.serverListener.Close()
	connId := 1
	for {
		conn, err := s.serverListener.Accept()
		if err != nil {
			s.log(err)
			break
		}
		ch := make(chan string)
		client := Client{connId, ch, 20, 20 + connId, 0, 0, 5}
		client.Listen(conn)
		s.Clients = append(s.Clients, client)
		go s.handleConn(conn, client)
		connId++
	}
	return nil

}

func (s *Server) handleConn(conn net.Conn, client Client) {
	done := make(chan string)
	fmt.Println("trying to handle connection")
CONNECTION:
	for {
		go s.handleStream(conn, client, done)
		println("connection started")
		for {
			select {
			case <-done:
				println("Closing connection")
				break CONNECTION
			}
		}
	}
	conn.Close()
}

func (s *Server) handleStream(conn net.Conn, client Client, done chan string) {
	//	defer close(client.Message)
	bufc := bufio.NewReader(conn)
	s.InitializeStream(conn, client)
	for {
		line, _, err := bufc.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "exit" {
			done <- "Stream Closed"
		}
		s.Broadcast(string(line))
	}
}

func (s *Server) InitializeStream(conn net.Conn, client Client) {
	rand.Seed(time.Now().Unix())
	x := rand.Intn(32-3)+3
	y := rand.Intn(32-3)+3
	conn.Write([]byte(fmt.Sprintf("pos,%d,%d,%d\n", client.Id, x, y)))
}

func (s *Server) Broadcast(line string) {
	for _, client := range s.Clients {
		client.Message <- string(line)
	}
}

func main() {
	s := NewServer()
	err := s.Serve()
	if err != nil {
		println(err)
	}

}
