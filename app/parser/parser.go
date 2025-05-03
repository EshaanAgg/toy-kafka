package parser

type Parser struct {
	bytes []byte
	idx   int
}

func NewParser(bytes []byte) *Parser {
	return &Parser{
		bytes: bytes,
		idx:   0,
	}
}

// Returns the n bytes starting from the current index
// and increments the index by n.
func (p *Parser) getNextBytes(n int) []byte {
	if p.idx+n > len(p.bytes) {
		return nil
	}
	b := p.bytes[p.idx : p.idx+n]
	p.idx += n
	return b
}
