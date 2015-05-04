package main

import (
	"net"
	"log"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

type Client struct {
	Id int
	Message chan string
}

type Server struct {
	Fighter1Port uint16
	Hostname string
	Logger *log.Logger
	Clients []Client
	clientListener net.Listener
	serverListener net.Listener
}

func NewClient(Id int) *Client {
	ch := make(chan string)
	return &Client{Id,ch}
}

func (client *Client) Listen(conn net.Conn){
	go func() {
		for line := range client.Message {
			str := strings.Split(strings.TrimSpace(string(line)),",");
			iid,_ := strconv.Atoi(str[0])
			x,_ := strconv.Atoi(str[1])
			y,_ := strconv.Atoi(str[2])
			fmt.Println(str)
			conn.Write([]byte(fmt.Sprintf("%d,%d,%d\n", iid, x, y)))
		}
	}()
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
		conn, err := s.serverListener.Accept();
		if err != nil {
			s.log(err)
			break
		}
		ch := make(chan string)
		client := Client{connId, ch}
		client.Listen(conn)
		s.Clients = append(s.Clients,client)
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
		println ("connection started")
		for {
			select {
			case <-done:
				println("Closing connection")
				break CONNECTION;
			}
		}
	}
	conn.Close()
}

func (s *Server) handleStream(conn net.Conn, client Client, done chan string) {
	defer close(client.Message)
	bufc := bufio.NewReader(conn)
	s.InitializeStream(conn, client)
	for {
		line, _, err := bufc.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "exit" {
			done<-"Stream Closed"
		}
		s.Broadcast(string(line))
	}
}

func (s *Server) InitializeStream(conn net.Conn, client Client) {
	conn.Write([]byte(fmt.Sprintf("%d,20,21\n", client.Id)))
}

func (s *Server) Broadcast(line string) {
	for _, client := range s.Clients {
		client.Message<-string(line)
	}
}

func main () {
	s := NewServer()
	err := s.Serve()
	if err != nil {
		println(err)
	}

}
