package response

import "encoding/binary"


type [T any] Response struct {
    CorrelationID datatypes.CorrelationID
    Body T
}

func NewResponse[T any](correlationID int32, body T) *Response[T] {
	return &Response[T]{
        CorrelationID: correlationID,
        Body: body,
    }
}

func (r *Response) Bytes() ([]byte, error){   
    var dataBytes bytes.Buffer
    
    // Add correlation ID to the buffer as header
    err := r.CorrelationID.Marshal(&dataBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal correlation ID: %w", err)
    }

    // Add the body to the buffer
    if err := Unmarshal(&r.Body, &dataBytes); err != nil {
        return nil, fmt.Errorf("failed to marshal response body: %w", err)
    }


    // Add the message length
    l := dataBytes.Len()
    bytes := binary.BigEndian.AppendUint32(nil, uint32(l))
    bytes = append(bytes, dataBytes.Bytes()...)

    return bytes
}
