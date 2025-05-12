package datatypes

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type UInt32 uint32
type VarInt int64
type VarUInt uint64
type VarLong = VarInt

func (i *Int8) Unmarshal(p *Parser) error {
	b := p.getNextBytes(1)
	if b == nil {
		return errors.New("int8: not enough bytes")
	}
	*i = Int8(b[0])
	return nil
}

func (i Int8) Marshal(b *bytes.Buffer) error {
	return b.WriteByte(byte(i))
}

func (i *Int16) Unmarshal(p *Parser) error {
	b := p.getNextBytes(2)
	if b == nil {
		return errors.New("int16: not enough bytes")
	}
	*i = Int16(binary.BigEndian.Uint16(b))
	return nil
}

func (i Int16) Marshal(b *bytes.Buffer) error {
	return binary.Write(b, binary.BigEndian, uint16(i))
}

func (i *Int32) Unmarshal(p *Parser) error {
	b := p.getNextBytes(4)
	if b == nil {
		return errors.New("int32: not enough bytes")
	}
	*i = Int32(binary.BigEndian.Uint32(b))
	return nil
}

func (i Int32) Marshal(b *bytes.Buffer) error {
	return binary.Write(b, binary.BigEndian, uint32(i))
}

func (i *Int64) Unmarshal(p *Parser) error {
	b := p.getNextBytes(8)
	if b == nil {
		return errors.New("int64: not enough bytes")
	}
	*i = Int64(binary.BigEndian.Uint64(b))
	return nil
}

func (i Int64) Marshal(b *bytes.Buffer) error {
	return binary.Write(b, binary.BigEndian, uint64(i))
}

func (i *UInt32) Unmarshal(p *Parser) error {
	b := p.getNextBytes(4)
	if b == nil {
		return errors.New("uint32: not enough bytes")
	}
	*i = UInt32(binary.BigEndian.Uint32(b))
	return nil
}

func (i UInt32) Marshal(b *bytes.Buffer) error {
	return binary.Write(b, binary.BigEndian, uint32(i))
}

func (i *VarInt) Unmarshal(p *Parser) error {
	v, n := binary.Varint(p.bytes[p.idx:])
	if n == 0 {
		return errors.New("varint: not enough bytes to decode varint")
	}
	if n < 0 {
		return errors.New("varint: varint overflow")
	}

	p.idx += n
	*i = VarInt(v)
	return nil
}

func (i VarInt) Marshal(b *bytes.Buffer) error {
	d := binary.AppendVarint([]byte{}, int64(i))
	_, err := b.Write(d)
	return err
}

func (i *VarUInt) Unmarshal(p *Parser) error {
	v, n := binary.Uvarint(p.bytes[p.idx:])
	if n == 0 {
		return errors.New("varuint: not enough bytes to decode varint")
	}
	if n < 0 {
		return errors.New("varuint: varint overflow")
	}

	p.idx += n
	*i = VarUInt(v)
	return nil
}

func (i VarUInt) Marshal(b *bytes.Buffer) error {
	d := binary.AppendUvarint([]byte{}, uint64(i))
	_, err := b.Write(d)
	return err
}
