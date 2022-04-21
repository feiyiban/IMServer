package server

import (
	"net"
)

type User struct {
	Name    string
	Address string

	C    chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Address: userAddr,
		C:       make(chan string),
		Conn:    conn,
	}

	go user.ListenMessage()

	return user
}

func (user *User) ListenMessage() {
	for {
		message := <-user.C
		user.Conn.Write([]byte(message + "\n"))
	}

}
