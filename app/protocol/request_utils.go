package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Decodable interface {
	Decode(r *Request)
}

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

func (r *Request) ReadInt8() int8 {
	val := int8(r.bytes[r.byteParseOffset])
	r.byteParseOffset++
	return val
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

func (r *Request) ReadInt64() int64 {
	val := binary.BigEndian.Uint64(r.bytes[r.byteParseOffset:])
	r.byteParseOffset += 8
	return int64(val)
}

func (r *Request) ReadVarInt() int {
	val := 0
	shift := 0

	for {
		b := r.bytes[r.byteParseOffset]
		r.byteParseOffset++

		val |= int(b&0x7F) << shift
		shift += 7

		if b&0x80 == 0 {
			break
		}
	}

	return val
}

func (r *Request) ReadUUID() string {
	bytes := r.bytes[r.byteParseOffset : r.byteParseOffset+16]
	r.byteParseOffset += 16
	return string(bytes)
}

func (r *Request) ReadCompactString() string {
	length := uint(r.ReadVarInt() - 1)

	bytes := r.bytes[r.byteParseOffset : r.byteParseOffset+length]
	r.byteParseOffset += length

	return string(bytes)
}

// DecodeInt32Array decodes an array of int32 values. This is a special case
// as we can't define additional methods to satify Decode interface on primitive types.
func DecodeInt32Array(r *Request) []int32 {
	count := r.ReadVarInt() - 1
	items := make([]int32, count)

	for i := range items {
		items[i] = r.ReadInt32()
	}

	return items
}

func (r *Request) ReadTaggedFields() {
	// Expect the tagged fields to be NULL
	if r.bytes[r.byteParseOffset] != 0 {
		fmt.Printf("Tagged fields are not NULL: %d\n", r.bytes[r.byteParseOffset])
		os.Exit(1)
	}

	r.byteParseOffset++
}
