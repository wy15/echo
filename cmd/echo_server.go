package main

import "echo/server"

func main() {
	key := make([]byte, 32)
	copy(key, "this is my key value!")
	server.UdpServe("localhost:8080", key)
	//server.TcpServe("localhost:8080", key)
}
