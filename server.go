package service

import "net"

type Server struct {
	IP   string
	Port string
}

func NewServer(ip string, port string) *Server {
	server := Server{
		IP:   ip,
		Port: port,
	}

	return &server
}

func (this *Server) Start() {
	// socker listen
	listener, err := net.Listen("tcp", fmt.sprintf("%s:%s", this.IP, this.Port))
}
