package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	"reflect"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      map[uint16]*MsgInfo
	msgID        map[reflect.Type]uint16
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
	msgRawMergedHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      uint16
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgID = make(map[reflect.Type]uint16)
	p.msgInfo = make(map[uint16]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msg proto.Message, id uint16) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatalf("protobuf message pointer required")
	}
	if _, ok := p.msgID[msgType]; ok {
		log.Fatalf("message %s is already registered", msgType)
	}
	if _, ok := p.msgInfo[id]; ok {
		log.Fatalf("message %v is already registered", id)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[id] = i
	p.msgID[msgType] = id
	return id
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatalf("message %s not registered", msgType)
	}
	if _, ok := p.msgInfo[id]; !ok {
		log.Fatalf("message id %v is already registered", id)
	}
	p.msgInfo[id].msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatalf("message %s not registered", msgType)
	}
	if _, ok := p.msgInfo[id]; !ok{
		log.Fatalf("message %v not registered", id)
	}
	p.msgInfo[id].msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(id uint16, msgRawHandler MsgHandler) {
	if _, ok := p.msgInfo[id]; !ok{
		log.Fatalf("message id %v not registered", id)
	}

	p.msgInfo[id].msgRawHandler = msgRawHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawMergedHandler(id uint16, msgRawMergedHandler MsgHandler) {
	if _, ok := p.msgInfo[id]; !ok{
		log.Fatalf("message id %v not registered", id)
	}

	p.msgInfo[id].msgRawMergedHandler = msgRawMergedHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		_, ok := p.msgInfo[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message id %v not registered", msgRaw.msgID)
		}
		i := p.msgInfo[msgRaw.msgID]
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
        if i.msgRawMergedHandler != nil {
			i.msgRawMergedHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// protobuf
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		return fmt.Errorf("message %s not registered", msgType)
	}
	if _, ok := p.msgInfo[id]; !ok{
		return fmt.Errorf("message id %v not registered", id)
	}
	i := p.msgInfo[id]
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRawHandler != nil {
		i.msgRawHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}

	// id
	var id uint16
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}
	if _, ok := p.msgInfo[id]; !ok {
		return nil, fmt.Errorf("message id %v not registered", id)
	}

	// msg
	i := p.msgInfo[id]
	if i.msgRawHandler != nil {
		return MsgRaw{id, data[2:]}, nil
	} else if i.msgRawMergedHandler != nil {
		return MsgRaw{id, data[0:]}, nil
	}else {
		msg := reflect.New(i.msgType.Elem()).Interface()
		return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
	}
}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)

	// id
	_id, ok := p.msgID[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := make([]byte, 2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, _id)
	} else {
		binary.BigEndian.PutUint16(id, _id)
	}

	// data
	data, err := proto.Marshal(msg.(proto.Message))
	return [][]byte{id, data}, err
}

// goroutine safe, return msg len and data as single data
func (p *Processor) MarshalWithMsgId(msg interface{}) ([]byte, error) {
	args, err := p.Marshal(msg)
	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// don't copy
	if len(args) == 1 {
		return args[0], nil
	}

	// merge the args
	msgMerge := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msgMerge[l:], args[i])
		l += len(args[i])
	}
	return msgMerge,err
}

// goroutine safe
func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}
