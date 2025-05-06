package datatypes

import "bytes"

type MarshableDataType interface {
	Marshal(b *bytes.Buffer) error
}

type UnmarshableDataType interface {
	Unmarshal(p *Parser) error
}
