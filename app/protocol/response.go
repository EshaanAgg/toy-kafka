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

func (r *Response) WriteInt16(vals ...int16) {
	for _, val := range vals {
		r.bytes = binary.BigEndian.AppendUint16(r.bytes, uint16(val))
	}
}

func (r *Response) WriteInt32(vals ...int32) {
	for _, val := range vals {
		r.bytes = binary.BigEndian.AppendUint32(r.bytes, uint32(val))
	}
}

func (r *Response) Send(conn *net.Conn) {
	msg := make([]byte, 0)

	msgLen := 4 + uint32(len(r.bytes))
	msg = binary.BigEndian.AppendUint32(msg, msgLen)
	msg = append(msg, r.bytes...)

	_, err := (*conn).Write(msg)
	if err != nil {
		fmt.Printf("Error sending response: %s\n", err.Error())
	}
}
