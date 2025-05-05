package request

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes/request/parser"
)

type HeaderFields struct {
	APIKey        int16   `kafka:"int16"`
	APIVersion    int16   `kafka:"int16"`
	CorrelationID int32   `kafka:"int32"`
	ClientID      *string `kafka:"nullable_string"`
}

type RequestHeader struct {
	// Fields used for parsing
	*parser.Parser

	Length int32
	*HeaderFields
}

func NewRequestHeader(buf []byte) (*RequestHeader, error) {
	r := &RequestHeader{
		Parser: parser.NewParser(buf),
	}

	bodyLen, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [messageLength]: %w", err)
	}
	r.Length = bodyLen

	var hf HeaderFields
	if err := r.AutoDecodeBody("HeaderFields", &hf, true); err != nil {
		return nil, err
	}
	r.HeaderFields = &hf

	return r, nil
}
