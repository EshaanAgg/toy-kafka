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

func (p *Parser) ReadVarInt() (uint, error) {
	var result uint
	var shift uint

	for {
		if p.idx >= len(p.bytes) {
			return 0, errors.New("ReadVarInt: not enough bytes")
		}
		b := p.bytes[p.idx]
		p.idx++
		result |= (uint(b) & 0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	return result, nil
}

func (p *Parser) ReadNullableString() (*string, error) {
	l, err := p.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("ReadNullableString [length]: %w", err)
	}

	if l == -1 {
		return nil, nil
	}

	content := p.getNextBytes(int(l))
	if content == nil {
		return nil, fmt.Errorf("ReadNullableString [content]: Not enough bytes for length %d", l)
	}

	s := string(content)
	return &s, nil
}

func (p *Parser) ReadCompactString() (string, error) {
	// The value of N + 1 is encoded as an UNSIGNED variable-length integer
	n, err := p.ReadVarInt()
	if err != nil {
		return "", fmt.Errorf("ReadCompactString [length]: %w", err)
	}
	// Type-cast to int and decrement by 1
	n = uint(n) - 1

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return "", fmt.Errorf("ReadCompactString [content]: Not enough bytes for length %d", n)
	}
	return string(strBytes), nil
}
