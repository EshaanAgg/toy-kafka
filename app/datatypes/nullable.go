package datatypes

import (
	"bytes"
	"errors"
	"fmt"
)

// Nullable is a genetic type that can be used to represent a nullable value.
// The type T should be a struct.
type Nullable[T any] struct {
	Present bool
	Value   *T
}

func NewNullable[T any]() Nullable[T] {
	return Nullable[T]{
		Present: false,
		Value:   nil,
	}
}

func (n *Nullable[T]) Unmarshal(p *Parser) error {
	// Read the first byte to check if the value is present
	presentByte := p.getNextBytes(1)
	if presentByte == nil {
		return errors.New("nullable: unable to read 1 byte")
	}
	if presentByte[0] == 0xFF {
		n.Present = false
		n.Value = nil
	} else {
		n.Present = true
		value := new(T)
		if err := Unmarshal(value, p); err != nil {
			return fmt.Errorf("nullable: unable to unmarshal value: %w", err)
		}
		n.Value = value
	}

	return nil
}

func (n Nullable[T]) Marshal(b *bytes.Buffer) error {
	if n.Present {
		if err := b.WriteByte(0x00); err != nil {
			return fmt.Errorf("nullable: unable to write present byte: %w", err)
		}
		if err := Marshal(n.Value, b); err != nil {
			return fmt.Errorf("nullable: unable to marshal value: %w", err)
		}
	} else {
		if err := b.WriteByte(0xFF); err != nil {
			return fmt.Errorf("nullable: unable to write null byte: %w", err)
		}
	}
	return nil
}
