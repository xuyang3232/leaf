package gate

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"net"
	"reflect"
)

type UDPGate struct {
	Processor     network.Processor
	UDPAddr      string
	PendingWriteNum int
}

type UDPRouteData struct{
	A    *UDPAgent
	Addr *net.UDPAddr
}

func (gate *UDPGate) Run(closeSig chan bool) {
	var udpServer *network.UDPServer
	if gate.UDPAddr != "" {
		udpServer = new(network.UDPServer)
		udpServer.Addr = gate.UDPAddr
		udpServer.PendingWriteNum = gate.PendingWriteNum
		udpServer.NewAgent = func(conn *network.UDPConn) network.Agent {
			a := &UDPAgent{conn: conn, gate: gate}
			return a
		}
	}

	if udpServer != nil {
		udpServer.Start()
	}
	<-closeSig

	if udpServer != nil {
		udpServer.Close()
	}
}

func (gate *UDPGate) OnDestroy() {}

type UDPAgent struct {
	conn     *network.UDPConn
	gate     *UDPGate
	userData interface{}
}

func (a *UDPAgent) Run() {
	for {
		data, addr, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			continue
		}
		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				continue
			}
			
			err = a.gate.Processor.Route(msg, &UDPRouteData{a, addr})
			if err != nil {
				log.Debug("route message error: %v", err)
				continue
			}
		}
	}
}
func (a *UDPAgent) OnClose() {
	
}
func (a *UDPAgent) WriteMsg(msg interface{}, addr * net.UDPAddr) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(addr,data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *UDPAgent) UserData() interface{} {
	return a.userData
}

func (a *UDPAgent) SetUserData(data interface{}) {
	a.userData = data
}