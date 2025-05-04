package parser

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func (p *Parser) ReadInt8() (int8, error) {
	b := p.getNextBytes(1)
	if b == nil {
		return 0, errors.New("ReadInt8: not enough bytes")
	}
	return int8(b[0]), nil
}

func (p *Parser) ReadInt16() (int16, error) {
	b := p.getNextBytes(2)
	if b == nil {
		return 0, errors.New("ReadInt16: not enough bytes")
	}

	// Can type-cast directly to int16 as the memory layout is the same
	return int16(binary.BigEndian.Uint16(b)), nil
}

func (p *Parser) ReadInt32() (int32, error) {
	b := p.getNextBytes(4)
	if b == nil {
		return 0, errors.New("ReadInt32: not enough bytes")
	}

	// Can type-cast directly to int32 as the memory layout is the same
	return int32(binary.BigEndian.Uint32(b)), nil
}

func (p *Parser) ReadInt64() (int64, error) {
	b := p.getNextBytes(8)
	if b == nil {
		return 0, errors.New("ReadInt64: not enough bytes")
	}

	// Can type-cast directly to int64 as the memory layout is the same
	return int64(binary.BigEndian.Uint64(b)), nil
}

func (p *Parser) ReadVarInt() (int64, error) {
	v, n := binary.Varint(p.bytes[p.idx:])
	if n == 0 {
		return 0, fmt.Errorf("ReadVarUInt: not enough bytes to decode varint")
	}
	if n < 0 {
		return 0, fmt.Errorf("ReadVarUInt: varint overflow")
	}

	p.idx += n
	return v, nil
}

func (p *Parser) ReadVarUInt() (uint64, error) {
	v, n := binary.Uvarint(p.bytes[p.idx:])
	if n == 0 {
		return 0, fmt.Errorf("ReadVarUInt: not enough bytes to decode varint")
	}
	if n < 0 {
		return 0, fmt.Errorf("ReadVarUInt: varint overflow")
	}

	p.idx += n
	return v, nil
}
