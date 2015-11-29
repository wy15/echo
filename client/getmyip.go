package main

import (
	"echo/libsodium"
	"flag"
	"log"
	"net"
	"time"
)

var key = flag.String("key", "this is my key value!", "key is 32 bytes")
var addr = flag.String("server", "localhost:8080", "like this localhost:8080")
var message = flag.String("message", "this is my home!", "some text")

func main() {
	flag.Parse()
	k := make([]byte, 32)
	copy(k, []byte(*key))
	r, err := tcpClient(*addr, k, []byte(*message))
	if err != nil {
		log.Printf("tcpclient error %v\n", err)
	} else {
		log.Printf("home ip is %s\n", r)
	}
}

func tcpClient(addr string, key, message []byte) ([]byte, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	ciphertext, err := libsodium.EncryptData(key, message)
	if err != nil {
		return nil, err
	}

	tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		return nil, err
	}
	defer tcpconn.Close()

	tcpconn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))
	_, err = tcpconn.Write(ciphertext)
	if err != nil {
		return nil, err
	}

	rtn := make([]byte, 21)
	tcpconn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))
	n, err := tcpconn.Read(rtn)
	if err != nil {
		return nil, err
	}

	return rtn[:n], nil
}
