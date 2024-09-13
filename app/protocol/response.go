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

func (r *Response) WriteVarInt(vals ...int) {
	for _, val := range vals {
		bytes := make([]byte, 0)

		for {
			toWrite := byte(val & 0x7F)
			val >>= 7
			if val != 0 {
				toWrite |= 0x80
				bytes = append(bytes, toWrite)
			} else {
				bytes = append(bytes, toWrite)
				break
			}
		}

		r.bytes = append(r.bytes, bytes...)
	}
}

func (r *Response) WriteBytes(vals ...byte) {
	r.bytes = append(r.bytes, vals...)
}

func (r *Response) Send(conn *net.Conn) {
	msg := make([]byte, 0)

	msgLen := uint32(len(r.bytes))
	msg = binary.BigEndian.AppendUint32(msg, msgLen)
	msg = append(msg, r.bytes...)

	_, err := (*conn).Write(msg)
	if err != nil {
		fmt.Printf("Error sending response: %s\n", err.Error())
	}
}
