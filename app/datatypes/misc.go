package datatypes

import (
	"bytes"
	"errors"
	"fmt"
)

type TaggedFields int
type UUID [16]byte
type CompactRecords = VarIntBytes

func (t *TaggedFields) Unmarshal(p *Parser) error {
	var tagLen VarUInt = 0
	if err := tagLen.Unmarshal(p); err != nil {
		return fmt.Errorf("tag_length: %w", err)
	}
	return nil
}

func (t TaggedFields) Marshal(b *bytes.Buffer) error {
	var tagLen VarUInt = 0
	if err := tagLen.Marshal(b); err != nil {
		return fmt.Errorf("tag_length: %w", err)
	}
	return nil
}

func (u *UUID) Unmarshal(p *Parser) error {
	b := p.getNextBytes(16)
	if b == nil {
		return errors.New("UUID: unable to read 16 bytes")
	}

	// Bytes are in big-endian order, so we need to reverse them to get the correct UUID
	for i := range 16 {
		u[i] = b[15-i]
	}
	return nil
}

func (u UUID) Marshal(b *bytes.Buffer) error {
	// Bytes are to be encoded in big-endian order
	for i := range 16 {
		if err := b.WriteByte(u[15-i]); err != nil {
			return fmt.Errorf("UUID: unable to write byte %d: %w", i, err)
		}
	}
	return nil
}
