package datatypes

import (
	"bytes"
	"fmt"
)

type TaggedFields int

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
