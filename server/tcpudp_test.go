package server_test

import (
	"bufio"
	"echo/netstring"
	"echo/server"
	"net"
	"testing"
	"time"
)

func TestTcpServe(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "this is my key value!")
	//go server.TcpServe("localhost:8080", key)

	tcpaddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("ResolveTCPAddr error : %v\n", err)
	}
	tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		t.Fatalf("DialTCP error : %v\n", err)
	}
	defer tcpconn.Close()
	//tcpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	plaintext := []byte("this is plaintext")
	ciphertext, err := server.EncryptData(key, plaintext, nil)
	if err != nil {
		t.Fatalf("EncryptData error : %v\n", err)
	}

	bufWriter := bufio.NewWriter(tcpconn)
	_, err = bufWriter.Write(netstring.Marshall(ciphertext))
	if err != nil {
		t.Fatalf("Write error : %v\n", err)
	}
	bufWriter.Flush()

	rdata := make([]byte, 50)
	rlen, err := tcpconn.Read(rdata)
	if err != nil {
		t.Fatalf("Read error : %v\n", err)
	} else {
		t.Logf("return data is %s\n", rdata[:rlen])

	}
}

func TestUdpServe(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "this is my key value!")

	udpaddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		t.Fatalf("ResolveUDPAddr error : %v\n", err)
	}

	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		t.Fatalf("DialUDP error : %v\n", err)
	}
	defer udpconn.Close()
	udpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))

	plaintext := []byte("a")
	ciphertext, err := server.EncryptData(key, plaintext, nil)
	if err != nil {
		t.Fatalf("EncryptData error : %v\n", err)
	}

	_, err = udpconn.Write(ciphertext)
	if err != nil {
		t.Fatalf("Write error : %v\n", err)
	}
}
