package datatypes

import (
	"bytes"
	"fmt"
)

type String string
type CompactString string
type VarIntString string

type NullableString struct {
	IsNull bool
	Value  string
}

func (s *String) Unmarshal(p *Parser) error {
	var n Int16
	if err := n.Unmarshal(p); err != nil {
		return fmt.Errorf("string.length: %w", err)
	}

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return fmt.Errorf("string.content: Not enough bytes for length %d", n)
	}
	*s = String(strBytes)
	return nil
}

func (s String) Marshal(b *bytes.Buffer) error {
	l := Int16(len(s))
	if err := l.Marshal(b); err != nil {
		return fmt.Errorf("string.length: %w", err)
	}
	if _, err := b.Write([]byte(s)); err != nil {
		return fmt.Errorf("string.content: %w", err)
	}
	return nil
}

func (s *CompactString) Unmarshal(p *Parser) error {
	// n + 1 is encoded as an unsigned variable-length integer
	var n VarUInt
	if err := n.Unmarshal(p); err != nil {
		return fmt.Errorf("compact_string.length: %w", err)
	}

	if n == 0 {
		*s = ""
		return nil
	}

	n -= 1
	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return fmt.Errorf("compact_string.content: Not enough bytes for length %d", n)
	}
	*s = CompactString(strBytes)
	return nil
}

func (s CompactString) Marshal(b *bytes.Buffer) error {
	l := VarUInt(len(s) + 1)
	if err := l.Marshal(b); err != nil {
		return fmt.Errorf("compact_string.length: %w", err)
	}
	if _, err := b.Write([]byte(s)); err != nil {
		return fmt.Errorf("compact_string.content: %w", err)
	}
	return nil
}

func (s *NullableString) Unmarshal(p *Parser) error {
	var n Int16
	if err := n.Unmarshal(p); err != nil {
		return fmt.Errorf("nullable_string.length: %w", err)
	}

	if n == -1 {
		s.IsNull = true
		return nil
	}

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return fmt.Errorf("nullable_string.content: Not enough bytes for length %d", n)
	}
	s.IsNull = false
	s.Value = string(strBytes)
	return nil
}

func (s NullableString) Marshal(b *bytes.Buffer) error {
	if s.IsNull {
		n := Int16(-1)
		if err := n.Marshal(b); err != nil {
			return fmt.Errorf("nullable_string.null: %w", err)
		}
		return nil
	}

	l := Int16(len(s.Value))
	if err := l.Marshal(b); err != nil {
		return fmt.Errorf("nullable_string.length: %w", err)
	}
	if _, err := b.Write([]byte(s.Value)); err != nil {
		return fmt.Errorf("nullable_string.content: %w", err)
	}
	return nil
}

func (s *VarIntString) Unmarshal(p *Parser) error {
	var n VarInt
	if err := n.Unmarshal(p); err != nil {
		return fmt.Errorf("varint_string.length: %w", err)
	}

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return fmt.Errorf("varint_string.content: Not enough bytes for length %d", n)
	}
	*s = VarIntString(strBytes)
	return nil
}

func (s VarIntString) Marshal(b *bytes.Buffer) error {
	l := VarInt(len(s))
	if err := l.Marshal(b); err != nil {
		return fmt.Errorf("varint_string.length: %w", err)
	}
	if _, err := b.Write([]byte(s)); err != nil {
		return fmt.Errorf("varint_string.content: %w", err)
	}
	return nil
}
