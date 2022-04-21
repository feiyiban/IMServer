package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, server)
	user.Online()

	// 监听用户时候活跃
	islive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err:", err)
				return
			}

			//提取用户的消息
			msg := string(buf[:n-1])

			//将用户消息广播
			user.DoMessage(msg)
			// 用户的任意消息
			islive <- true

		}
	}()

	for {
		select {
		case <-islive:
			//当前用户是活跃
			//不做任何处理
		case <-time.After(time.Second * 10):
			user.SendMsg("你被提了")
			// 销毁用的资源
			close(user.C)
			// 关闭连接
			conn.Close()

			return
		}
	}
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
