package datatypes

import (
	"bytes"
	"fmt"
)

type CompactNullableBytes []byte

// First the length N+1 is given as an UNSIGNED_VARINT.Then N bytes follow. A null object is represented with a length of 0.
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
