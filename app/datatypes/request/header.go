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

const MIN_SUPPORTED_API_VERSION = 0
const MAX_SUPPORTED_API_VERSION = 4

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

	apiVersion, err := r.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [apiVersion]: %w", err)
	}
	if apiVersion < MIN_SUPPORTED_API_VERSION || apiVersion > MAX_SUPPORTED_API_VERSION {
		return nil, fmt.Errorf("NewRequestHeader [apiVersion]: %w", fmt.Errorf("unsupported api version: %d", apiVersion))
	}

	correlationID, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [correlationID]: %w", err)
	}

	clientID, err := r.ReadNullableString()
	if err != nil {
		return nil, fmt.Errorf("NewRequestHeader [clientID]: %w", err)
	}

	return &RequestHeader{
		APIKey:        apiKey,
		APIVersion:    apiVersion,
		CorrelationID: correlationID,
		ClientID:      clientID,
	}, nil
}
