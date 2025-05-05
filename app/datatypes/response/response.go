package response

import "encoding/binary"

type Response struct {
	body []byte
}

func NewResponse(correlationID int32) *Response {
	r := &Response{}
	// Add correlation ID to the response by default
	r.WriteInt32(correlationID)

	return r
}

// Creates a new response with the given correlation ID and body.
// It automatically encodes the body using the AutoEncodeBody method.
// The base response struct is always assumed to NOT be an inline struct,
// and thus would end with empty tagged fields.
func NewResponseWithBody(correlationID int32, bodyPtr any) (*Response, error) {
	r := NewResponse(correlationID)
	if err := r.AutoEncodeBody(bodyPtr, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Response) Bytes() []byte {
	var bytes []byte

	bytes = binary.BigEndian.AppendUint32(bytes, uint32(len(r.body))) // Length of the response
	bytes = append(bytes, r.body...)                                  // Data bytes

	return bytes
}
