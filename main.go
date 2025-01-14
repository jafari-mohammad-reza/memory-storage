package main

import (
	"log"
)

type Config struct {
	ServerListenAddr string
}

func main() {
	server := NewServer(Config{
		ServerListenAddr: "8001",
	})
	log.Fatal(server.Start())
}
