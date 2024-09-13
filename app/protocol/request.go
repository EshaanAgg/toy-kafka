package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Request struct {
	bytes           []byte
	byteParseOffset uint

	Conn   *net.Conn
	Length int32
	Header *RequestHeader
}

// Utility functions to create request and parse it
func NewRequest(bytes []byte, conn *net.Conn) *Request {
	req := Request{
		bytes:           bytes,
		byteParseOffset: 0,
		Conn:            conn,
	}

	req.Length = req.ReadInt32()
	req.Header = ParseRequestHeader(&req)

	return &req
}

func (r *Request) ReadInt16() int16 {
	val := binary.BigEndian.Uint16(r.bytes[r.byteParseOffset:])
	r.byteParseOffset += 2
	return int16(val)
}

func (r *Request) ReadInt32() int32 {
	val := binary.BigEndian.Uint32(r.bytes[r.byteParseOffset:])
	r.byteParseOffset += 4
	return int32(val)
}

// Dispatch different handlers for different APIs
const API_VERSIONS_KEY = 18

type ResponseBody interface {
	Encode(r *Response)
}

func (r *Request) Handle() {
	var body ResponseBody

	switch r.Header.APIKey {
	case API_VERSIONS_KEY:
		body = GetAPIVersionBody(r)

	default:
		fmt.Printf("No handler has been implemented for the request API Key %d", r.Header.APIKey)
		return
	}

	resp := NewResponse()

	resp.WriteInt32(r.Header.CorrelationID) // Header -> CorrelationID
	body.Encode(resp)                       // Body Bytes
	resp.Send(r.Conn)
}
