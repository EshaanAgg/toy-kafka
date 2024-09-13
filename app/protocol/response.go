package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Response struct {
	bytes []byte
}

func NewResponse() *Response {
	return &Response{
		bytes: []byte{},
	}
}

func (r *Response) WriteInt16(val int16) {
	binary.BigEndian.AppendUint16(r.bytes, uint16(val))
}

func (r *Response) WriteInt32(val int32) {
	binary.BigEndian.AppendUint32(r.bytes, uint32(val))
}

func (r *Response) Send(conn *net.Conn) {
	_, err := (*conn).Write(r.bytes)
	if err != nil {
		fmt.Printf("Error writing: %s\n", err.Error())
	}
}
