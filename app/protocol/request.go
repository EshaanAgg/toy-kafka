package protocol

import "encoding/binary"

type Request struct {
	bytes           []byte
	byteParseOffset uint

	Length int32
	Header *RequestHeader
}

func NewRequest(bytes []byte) *Request {
	req := Request{
		bytes:           bytes,
		byteParseOffset: 0,
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
