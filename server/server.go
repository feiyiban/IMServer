package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int64
	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int64) *Server {
	srv := Server{
		IP:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return &srv
}

func (server *Server) Start() {
	// socker listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	fmt.Printf("server listen on %s:%d\n", server.IP, server.Port)
	defer listener.Close()
	go server.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}

		go server.handler(conn)

	}
}

func (server *Server) handler(conn net.Conn) {
	// 用户上线, 加入到OnLineMap
	user := NewUser(conn)
	server.mapLock.Lock()
	server.OnLineMap[user.Name] = user
	server.mapLock.Unlock()

	server.Broadcast(user, "上线啦!")

	// select {}
}

func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, cli := range server.OnLineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}
