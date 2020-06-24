package network

import (
	"github.com/name5566/leaf/log"
	"net"
)


type UDPServer struct {
	Addr            string
	NewAgent        func(conn *UDPConn) Agent
	PendingWriteNum int
	conn            *net.UDPConn
}

func (server *UDPServer) Start() {
	server.init()
	go server.run()
}

func (server *UDPServer) init() {
	udpAddr, err := net.ResolveUDPAddr("udp4", server.Addr)
	if err != nil{
		log.Fatalf("%v", err)
	}
	server.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		log.Infof("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	
	if server.NewAgent == nil {
		log.Fatalf("NewAgent must not be nil")
	}
}

func (server *UDPServer) run() {
	udpConn := newUDPConn(server.conn, server.PendingWriteNum)
	agent := server.NewAgent(udpConn)
	for {
		agent.Run()
	}
}

func (server *UDPServer) Close(){
	server.conn.Close()
}
