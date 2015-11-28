package main

import (
	"echo/server"
	"flag"
)

var key = flag.String("key", "this is my key value!", "key is 32 bytes")
var addr = flag.String("server", "localhost:8080", "like this localhost:8080")

func main() {
	flag.Parse()
	k := make([]byte, 32)
	copy(k, []byte(*key))
	go server.UdpServe(*addr, k)
	server.TcpServe(*addr, k)
}
