package tcpudp

import (
	"encrypt"
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

func udpServer(addr string) (*net.UDPConn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return net.ListenUDP("udp", udpaddr)
}

func TcpServe(addr string, encryptKey []byte) error {
	tcplistener, err := tcpListener(addr)
	if err != nil {
		return err
	}

	for {
		tcpconn, err := tcplistener.AcceptTCP()
		if err != nil {
			log.Printf("AcceptTCP error:%v", err)
			continue
			//return err
		}

		go handleTCPConn(tcpconn, encryptKey)
	}
}

func handleTCPConn(tcpconn *net.TCPConn, encryptKey []byte) {
	defer tcpconn.Close()
	tcpconn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	receiveData := make([]byte, 100, 1024*1024)
	var receiveDatalen int = 0
	for {
		datalen, err := tcpconn.Read(receiveData[receiveDatalen:])
		if err != nil {
			log.Printf("TCPConn Read error:%v", err)
			return
		}
		receiveDatalen += datalen
		if datalen == 0 {
			break
		}
	}

	if receiveDatalen == 0 {
		return
	}

	plaintext, err := encrypt.DecryptData(encryptKey, receiveData, nil)
	if err != nil {
		log.Printf("DecryptData error:%v", err)
		return
	}

	_, err = tcpconn.Write([]byte(homeip))
	if err != nil {
		log.Printf("tcpconn error:%v", err)
	}

}
