package main

import "echo/server"
import "flag"

var keyvalue = flag.String("key","this is my key value!", "key used to encrypt data doesn't longer than 32")
var addr = flag.String("addr", "localhost:8080", "Server's listen address,host:port")
func main() {
	flag.Parse()

	key := make([]byte, 32)

	copy(key, []byte(*keyvalue))
	go server.UdpServe(*addr, key)
	server.TcpServe(*addr, key)
}
