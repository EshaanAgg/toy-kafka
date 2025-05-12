package datatypes

import (
	"bytes"
	"fmt"
)

type CompactNullableBytes []byte
type VarIntBytes []byte

func (c *CompactNullableBytes) Unmarshal(p *Parser) error {
	var length VarUInt
	if err := length.Unmarshal(p); err != nil {
		return fmt.Errorf("compact_nullable_bytes.length: %w", err)
	}

	if length == 0 {
		*c = nil
		return nil
	}

	length--
	bytes := p.getNextBytes(int(length))
	if bytes == nil {
		return fmt.Errorf("compact_nullable_bytes: unable to read %d bytes", length)
	}
	*c = bytes
	return nil
}

func (c CompactNullableBytes) Marshal(b *bytes.Buffer) error {
	if c == nil {
		return VarUInt(0).Marshal(b)
	}
	if err := VarUInt(len(c) + 1).Marshal(b); err != nil {
		return fmt.Errorf("compact_nullable_bytes.length: %w", err)
	}
	if _, err := b.Write(c); err != nil {
		return fmt.Errorf("compact_nullable_bytes: unable to write %d bytes: %w", len(c), err)
	}
	return nil
}

func (v *VarIntBytes) Unmarshal(p *Parser) error {
	var l VarInt
	if err := l.Unmarshal(p); err != nil {
		return fmt.Errorf("varint_bytes.length: %w", err)
	}

	// Null bytes
	if l < 0 {
		*v = nil
		return nil
	}

	// Empty bytes
	if l == 0 {
		*v = []byte{}
		return nil
	}

	// Non-empty bytes
	bytes := p.getNextBytes(int(l))
	if bytes == nil {
		return fmt.Errorf("varint_bytes: unable to read %d bytes", l)
	}
	*v = bytes

	return nil
}

func (v VarIntBytes) Marshal(b *bytes.Buffer) error {
	if v == nil {
		return VarInt(-1).Marshal(b)
	}

	if err := VarInt(len(v)).Marshal(b); err != nil {
		return fmt.Errorf("varint_bytes.length: %w", err)
	}
	if _, err := b.Write(v); err != nil {
		return fmt.Errorf("varint_bytes: unable to write %d bytes: %w", len(v), err)
	}

	return nil
}
