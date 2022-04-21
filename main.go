package main

import (
	"github.com/feiyiban/IMServer/server"
)

func main() {
	srv := server.NewServer("127.0.0.1", 8888)
	srv.Start()
}
