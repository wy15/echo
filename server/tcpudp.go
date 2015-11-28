package server

import (
	"echo/libsodium"
	"log"
	"net"
	"time"
)

var homeip string = ":0"

func tcpListener(addr string) (*net.TCPListener, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", tcpaddr)
}

func TcpServe(addr string, encryptKey []byte) error {
	tcplistener, err := tcpListener(addr)
	if err != nil {
		return err
	}

	for {
		tcpconn, err := tcplistener.AcceptTCP()
		if err != nil {
			log.Printf("AcceptTCP error : %v", err)
			continue
		}

		go handleTCPConn(tcpconn, encryptKey)
	}
}

func handleTCPConn(tcpconn *net.TCPConn, encryptKey []byte) {
	defer tcpconn.Close()
	tcpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	receiveData := make([]byte, 50)
	receiveDatalen, err := tcpconn.Read(receiveData)
	if err != nil {
		log.Printf("TCPConn Read error : %v", err)
		return
	}

	if receiveDatalen == 0 {
		return
	}

	_, err = libsodium.DecryptData(encryptKey, receiveData[:receiveDatalen])
	if err != nil {
		log.Printf("DecryptData error : %v", err)
		return
	}

	_, err = tcpconn.Write([]byte(homeip))
	if err != nil {
		log.Printf("tcpconn error : %v", err)
	}
}

func UdpServe(addr string, key []byte) error {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpconn, err := net.ListenUDP("udp", udpaddr)
	for {
		if err != nil {
			log.Printf("udpServer error : %v", err)
			continue
		}
		handleUDPConn(udpconn, key)
	}
}

func handleUDPConn(udpconn *net.UDPConn, key []byte) {
	receiveData := make([]byte, 50)
	receiveDatalen, addr, err := udpconn.ReadFrom(receiveData)

	if err != nil {
		log.Printf("udp readfrom error : %v", err)
		return
	}

	_, err = libsodium.DecryptData(key, receiveData[:receiveDatalen])
	if err != nil {
		log.Printf("DecryptData error : %v", err)
		return
	}

	homeip = addr.String()
}
