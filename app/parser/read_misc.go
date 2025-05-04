package parser

import (
	"errors"
	"fmt"
)

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

func (p *Parser) ReadUUID() ([]byte, error) {
	b := p.getNextBytes(16)
	if b == nil {
		return nil, errors.New("ReadUUID: not enough bytes")
	}
	return b, nil
}
