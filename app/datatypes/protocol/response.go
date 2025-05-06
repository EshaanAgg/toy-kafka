package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

type Response[T any] struct {
	CorrelationID datatypes.Int32
	Body          *T
}

func NewResponse[T any](correlationID datatypes.Int32, body *T) *Response[T] {
	return &Response[T]{
		CorrelationID: correlationID,
		Body:          body,
	}
}

func (r *Response[T]) Bytes() ([]byte, error) {
	var dataBytes bytes.Buffer

	// Add correlation ID to the buffer as header
	err := r.CorrelationID.Marshal(&dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal correlation ID: %w", err)
	}

	// Add the body to the buffer
	if err := datatypes.Marshal(*r.Body, &dataBytes); err != nil {
		return nil, fmt.Errorf("failed to marshal response body: %w", err)
	}

	// Add the message length
	l := dataBytes.Len()
	bytes := binary.BigEndian.AppendUint32(nil, uint32(l))
	bytes = append(bytes, dataBytes.Bytes()...)

	return bytes, nil
}
