package protocol

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

type Request interface {
	Handle() ([]byte, error)
}

type HeaderFields struct {
	APIKey        datatypes.Int16
	APIVersion    datatypes.Int16
	CorrelationID datatypes.Int32
	ClientID      datatypes.NullableString
}

type RequestHeader struct {
	P      *datatypes.Parser
	Length int32
	*HeaderFields
}

func NewRequestHeader(buf []byte) (*RequestHeader, error) {
	r := &RequestHeader{
		P: datatypes.NewParser(buf),
	}

	var bodyLen datatypes.Int32
	if err := bodyLen.Unmarshal(r.P); err != nil {
		return nil, fmt.Errorf("error in unmarshalling messageLength: %w", err)
	}
	r.Length = int32(bodyLen)

	var hf HeaderFields
	if err := datatypes.Unmarshal(&hf, r.P); err != nil {
		return nil, fmt.Errorf("error in unmarshalling headerFields: %w", err)
	}
	r.HeaderFields = &hf

	return r, nil
}
