package network

import (
	"github.com/name5566/leaf/log"
	"net"
	"time"
)

type UDPClient struct {
	ConnectInterval time.Duration
	Addr            string
	NewAgent        func(conn *UDPConn) Agent
	PendingWriteNum int
}

func (client *UDPClient) Start() {
	client.init()
	go client.connect()
}

func (client *UDPClient) init() {
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 100
		log.Release("invalid PendingWriteNum, reset to %v", client.PendingWriteNum)
	}
	
	if client.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
}

func (client *UDPClient) dial() *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp4", client.Addr)
	if err != nil{
		return  nil
	}
	for {
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err == nil{
			return conn
		}

		log.Release("connect to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *UDPClient) connect() {
	conn := client.dial()
	if conn == nil {
		return
	}
	udpConn := newUDPConn(conn, client.PendingWriteNum)
	agent := client.NewAgent(udpConn)
	agent.Run()

	// cleanup
	conn.Close()
}

func (client *UDPClient) Close() {

}
