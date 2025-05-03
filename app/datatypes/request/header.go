package request

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/parser"
)

type RequestHeader struct {
	// Fields used for parsing
	*parser.Parser

	// Length of the request body in bytes
	Length int32

	// Header fields
	APIKey        int16
	APIVersion    int16
	CorrelationID int32
	ClientID      *string
}

func NewRequestHeader(buf []byte) (*RequestHeader, error) {
	r := &RequestHeader{
		Parser: parser.NewParser(buf),
	}

	bodyLen, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("NewRequest [messageLength]: %w", err)
	}
	r.Length = bodyLen

	apiKey, err := r.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [apiKey]: %w", err)
	}
	r.APIKey = apiKey

	apiVersion, err := r.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [apiVersion]: %w", err)
	}
	r.APIVersion = apiVersion

	correlationID, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [correlationID]: %w", err)
	}
	r.CorrelationID = correlationID

	clientID, err := r.ReadNullableString()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [clientID]: %w", err)
	}
	r.ClientID = clientID

	if err := r.ReadZeroTaggedFieldArray(); err != nil {
		return nil, fmt.Errorf("NewRequestHeader [taggedFields]: %w", err)
	}

	return r, nil
}
