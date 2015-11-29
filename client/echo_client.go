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
	for {
		err := udpClient(*addr, k, []byte(*message))
		if err != nil {
			log.Printf("udp client error %v\n", err)
		}

		time.Sleep(time.Duration(10) * time.Minute)
	}
}

func udpClient(addr string, key, message []byte) error {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	ciphertext, err := libsodium.EncryptData(key, message)
	if err != nil {
		return err
	}

	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return err
	}
	defer udpconn.Close()

	_, err = udpconn.Write(ciphertext)
	if err != nil {
		return err
	}

	return nil
}
