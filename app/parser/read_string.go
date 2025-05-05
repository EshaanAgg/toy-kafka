package parser

import "fmt"

func (p *Parser) ReadNullableString() (*string, error) {
	l, err := p.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("ReadNullableString [length]: %w", err)
	}

	if l == -1 {
		return nil, nil
	}

	content := p.getNextBytes(int(l))
	if content == nil {
		return nil, fmt.Errorf("ReadNullableString [content]: Not enough bytes for length %d", l)
	}

	v := string(content)
	return &v, nil
}

func (p *Parser) ReadCompactString() (string, error) {
	// The value of N + 1 is encoded as an unsigned variable-length integer
	n, err := p.ReadVarUInt()
	if err != nil {
		return "", fmt.Errorf("ReadCompactString [length]: %w", err)
	}
	if n == 0 {
		return "", nil
	}

	n -= 1
	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return "", fmt.Errorf("ReadCompactString [content]: Not enough bytes for length %d", n)
	}
	return string(strBytes), nil
}

func (p *Parser) ReadString() (string, error) {
	// The value of N is encoded as an INT16
	n, err := p.ReadInt16()
	if err != nil {
		return "", fmt.Errorf("ReadString [length]: %w", err)
	}

	strBytes := p.getNextBytes(int(n))
	if strBytes == nil {
		return "", fmt.Errorf("ReadString [content]: Not enough bytes for length %d", n)
	}
	return string(strBytes), nil
}
