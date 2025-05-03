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
	// The value of N + 1 is encoded as an unsigned variable-length integer
	n, err := p.ReadVarUInt()
	if err != nil {
		return "", fmt.Errorf("ReadCompactString [length]: %w", err)
	}
	if n == 0 {
		return "", nil
	}

	n -= 1
	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return "", fmt.Errorf("ReadCompactString [content]: Not enough bytes for length %d", n)
	}
	return string(strBytes), nil
}

func (p *Parser) ReadString() (string, error) {
	// The value of N + 1 is encoded as an UNSIGNED variable-length integer
	n, err := p.ReadInt16()
	if err != nil {
		return "", fmt.Errorf("ReadString [length]: %w", err)
	}

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return "", fmt.Errorf("ReadString [content]: Not enough bytes for length %d", n)
	}
	return string(strBytes), nil
}

func (p *Parser) ReadZeroTaggedFieldArray() error {
	// The array is encoded as an compact array
	n, err := p.ReadVarUInt()
	if err != nil {
		return fmt.Errorf("ReadZeroTaggedFieldArray [length]: %w", err)
	}
	if n != 0 {
		return fmt.Errorf("ReadZeroTaggedFieldArray [length]: expected 0, got %d", n)
	}
	return nil
}
