package network

import (
	"github.com/name5566/leaf/log"
	"net"
	"sync"
)

type UDPWriteData struct{
	data []byte
	userAddr *net.UDPAddr
}

type UDPConn struct {
	sync.Mutex
	conn      *net.UDPConn
	writeChan chan *UDPWriteData
	closeFlag bool
}

func newUDPConn(conn *net.UDPConn, pendingWriteNum int) *UDPConn {
	udpConn := new(UDPConn)
	udpConn.conn = conn
	udpConn.writeChan = make(chan *UDPWriteData, pendingWriteNum)

	go func() {
		for b := range udpConn.writeChan {
			if b == nil {
				break
			}
			_, _,err := udpConn.conn.WriteMsgUDP(b.data, nil, b.userAddr)
			if err != nil {
				continue
			}
		}
		conn.Close()
		udpConn.Lock()
		udpConn.closeFlag = true
		udpConn.Unlock()
	}()

	return udpConn
}

func (udpConn *UDPConn) doDestroy() {
	udpConn.conn.Close()

	if !udpConn.closeFlag {
		close(udpConn.writeChan)
		udpConn.closeFlag = true
	}
}

func (udpConn *UDPConn) Close() {
	udpConn.Lock()
	defer udpConn.Unlock()
	if udpConn.closeFlag {
		return
	}

	udpConn.doWrite(nil)
	udpConn.closeFlag = true
}

func (udpConn *UDPConn) doWrite(data *UDPWriteData) {
	if len(udpConn.writeChan) == cap(udpConn.writeChan) {
		log.Debug("close conn: channel full")
		udpConn.doDestroy()
		return
	}

	udpConn.writeChan <- data
}

func (udpConn *UDPConn) WriteMsg(addr *net.UDPAddr, args ...[]byte) error {
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	msg := make([]byte, uint32(msgLen))
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	udpConn.doWrite(&UDPWriteData{
		data:msg,
		userAddr:addr,
	})
	
	return  nil
}

func (udpConn *UDPConn) ReadMsg() ([]byte, *net.UDPAddr,error) {
	var buf [1024]byte
	n, addr, err := udpConn.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil, nil, err
	}
	return buf[0:n], addr, err
}

