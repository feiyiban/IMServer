package server

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Address string

	C    chan string
	Conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Address: userAddr,
		C:       make(chan string),
		Conn:    conn,
		server:  server,
	}

	go user.ListenMessage()

	return user
}

// 用户上线业务
func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnLineMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "上线啦!")
}

// 用户下线业务
func (user *User) Offline() {
	user.server.mapLock.Lock()
	user.server.OnLineMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "下线啦!")
}

// 发送消息给客户端
func (user *User) SendMsg(msg string) {
	user.Conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		user.server.mapLock.Lock()
		for _, user := range user.server.OnLineMap {
			onlineMsg := "[" + user.Address + "]" + user.Name + "online...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := user.server.OnLineMap[newName]
		if ok {
			user.SendMsg("当前用户名被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnLineMap, user.Name)
			user.server.OnLineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("您已经更新用户名:" + user.Name + "\n")
		}
	} else {
		user.server.Broadcast(user, msg)
	}

}

func (user *User) ListenMessage() {
	for {
		message := <-user.C
		user.Conn.Write([]byte(message + "\n"))
	}
}
