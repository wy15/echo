package main

import "echo/server"

func main() {
	key := make([]byte, 32)

	copy(key, "this is a key")
	server.TcpServe("localhost:8080", key)
}
