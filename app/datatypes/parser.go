package datatypes

type Parser struct {
	bytes []byte
	idx   int
	l     int
}

func NewParser(bytes []byte) *Parser {
	return &Parser{
		bytes: bytes,
		idx:   0,
		l:     len(bytes),
	}
}

// Returns the n bytes starting from the current index
// and increments the index by n.
func (p *Parser) getNextBytes(n int) []byte {
	if p.idx == p.l || p.idx+n > p.l {
		return nil
	}
	b := p.bytes[p.idx : p.idx+n]
	p.idx += n
	return b
}
