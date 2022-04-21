package server

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int64
}

func NewServer(ip string, port int64) *Server {
	srv := Server{
		IP:   ip,
		Port: port,
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

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}

		fmt.Printf("client %s connected\n", conn.RemoteAddr().String())
		go server.handler(&conn)

	}
}

func (server *Server) handler(conn *net.Conn) {
	// Handle msg
}
