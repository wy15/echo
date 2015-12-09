package main

import (
	"bufio"
	"bytes"
	"echo/libsodium"
	"echo/netstring"
	"flag"
	"io"
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

	tcpconn.SetWriteDeadline(time.Now().Add(time.Duration(20) * time.Second))
	_, err = tcpconn.Write(netstring.Marshall(ciphertext))
	if err != nil {
		return nil, err
	}

	tcpconn.SetReadDeadline(time.Now().Add(time.Duration(20) * time.Second))
	bufReader := bufio.NewReader(tcpconn)
	var buf bytes.Buffer
	var ip []byte
	for {
		rData, err := bufReader.ReadBytes(',')
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			buf.Write(rData)
			continue
		}
		buf.Write(rData)
		ip, err = netstring.Unmarshall(buf.Bytes())
		if err != nil {
			if err == netstring.ErrNsLenNotEqaulOrgLen {
				continue
			}
			return nil, err
		}
		break
	}
	return ip, nil
}
