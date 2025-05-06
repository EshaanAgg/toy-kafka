package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

type Response[H any, B any] struct {
	Header *H
	Body   *B
}

func NewResponse[H any, B any](header *H, body *B) *Response[H, B] {
	return &Response[H, B]{
		Header: header,
		Body:   body,
	}
}

func (r *Response[T, B]) Bytes() ([]byte, error) {
	var dataBytes bytes.Buffer

	// Add the header to the buffer
	err := datatypes.Marshal(*r.Header, &dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the header: %w", err)
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
