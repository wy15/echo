package server

import (
	"bufio"
	"bytes"
	"echo/bpool"
	"echo/libsodium"
	"echo/netstring"
	"io"
	"log"
	"net"
	"os"
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
	log.SetOutput(os.Stdout)
	if err != nil {
		return err
	}
	bpool := bpool.NewSizedBufferPool(20, 4096)

	for {
		tcpconn, err := tcplistener.AcceptTCP()
		if err != nil {
			log.Printf("AcceptTCP error : %v\n", err)
			continue
		}
		buffer := bpool.Get()

		go handleTCPConn(tcpconn, encryptKey, buffer)
	}
}

func handleTCPConn(tcpconn *net.TCPConn, encryptKey []byte, buffer *bytes.Buffer) {
	defer tcpconn.Close()
	//receiveData := make([]byte, 50)
	var receiveData []byte
	//tcpconn need read all data in 20 second ,otherwise Timeout() will be true
	//tcpconn.SetReadDeadline(time.Now().Add(time.Duration(20) * time.Second))
	//receiveDatalen, err := tcpconn.Read(receiveData)
	bufReader := bufio.NewReader(tcpconn)
	for {
		rData, err := bufReader.ReadBytes([]byte(",")[0])
		if err != nil {
			if err != io.EOF {
				log.Printf("TCPConn Read error : %v\n", err)
				return
			}
			buffer.Write(rData)
			continue
		}

		buffer.Write(rData)

		unmarshallData, err := netstring.Unmarshall(buffer.Bytes())
		if err == netstring.ErrNsLenNotEqaulOrgLen {
			continue
		}
		if err != nil {
			log.Printf("netstring unmarshall error : %v\n", err)
			return
		}
		receiveData = unmarshallData

		break
	}

	_, err := libsodium.DecryptData(encryptKey, receiveData)
	if err != nil {
		log.Printf("DecryptData error : %v\n", err)
		return
	}

	tcpconn.SetWriteDeadline(time.Now().Add(time.Duration(20) * time.Second))
	_, err = tcpconn.Write(netstring.Marshall([]byte(homeip)))
	if err != nil {
		log.Printf("tcpconn error : %v\n", err)
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
			log.Printf("udpServer error : %v\n", err)
			continue
		}
		handleUDPConn(udpconn, key)
	}
}

func handleUDPConn(udpconn *net.UDPConn, key []byte) {
	receiveData := make([]byte, 50)
	receiveDatalen, addr, err := udpconn.ReadFrom(receiveData)

	if err != nil {
		log.Printf("udp readfrom error : %v\n", err)
		return
	}

	_, err = libsodium.DecryptData(key, receiveData[:receiveDatalen])
	if err != nil {
		log.Printf("DecryptData error : %v\n", err)
		return
	}

	homeip = addr.String()
}
