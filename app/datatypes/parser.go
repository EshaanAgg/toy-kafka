package datatypes

import "fmt"

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
	if n < 0 {
		panic("getNextBytes: n must be non-negative")
	}
	if n == 0 {
		return make([]byte, 0)
	}

	if p.idx+n-1 >= p.l {
		return nil
	}
	b := p.bytes[p.idx : p.idx+n]
	p.idx += n
	return b
}

func (p *Parser) IsAtEnd() bool {
	return p.idx >= p.l
}

func (p *Parser) Debug() {
	fmt.Printf("Parser [idx=%d, l=%d]: ", p.idx, p.l)
	end := min(p.idx+10, p.l)
	fmt.Printf("bytes=%v\n", p.bytes[p.idx:end])
}
