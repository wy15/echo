package server_test

import (
	"echo/server"
	"net"
	"testing"
	"time"
)

func TestTcpServe(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "this is a key ddddldlldd *** d000 d877668786566")
	//go server.TcpServe("localhost:8080", key)

	tcpaddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		t.Errorf("ResolveTCPAddr error:%v", err)
	}
	tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		t.Errorf("DialTCP error:%v", err)
	}
	defer tcpconn.Close()
	tcpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	plaintext := []byte("this is plaintext")
	ciphertext, err := server.EncryptData(key, plaintext, nil)
	if err != nil {
		t.Errorf("EncryptData error:%v", err)
	}

	_, err = tcpconn.Write(ciphertext)
	if err != nil {
		t.Errorf("Write error:%v", err)
	}

	rdata := make([]byte, 50)
	rlen, err := tcpconn.Read(rdata)
	if err != nil {
		t.Errorf("Read error:%v", err)
	} else {
		t.Logf("return data is %s", rdata[:rlen])

	}
}

func TestUdpServe(t *testing.T) {
	key := make([]byte, 32)
	copy(key, "this is a key ddddldlldd *** d000 d877668786566")

	udpaddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		t.Errorf("ResolveUDPAddr error:%v", err)
	}

	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		t.Errorf("DialUDP error:%v", err)
	}
	defer udpconn.Close()
	udpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))

	plaintext := []byte("this is plaintext")
	ciphertext, err := server.EncryptData(key, plaintext, nil)
	if err != nil {
		t.Errorf("EncryptData error:%v", err)
	}

	_, err = udpconn.Write(ciphertext)
	if err != nil {
		t.Errorf("Write error:%v", err)
	}
}
