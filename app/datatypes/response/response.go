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

func (r *Response) Bytes() []byte {
	var bytes []byte

	bytes = binary.BigEndian.AppendUint32(bytes, uint32(len(r.body))) // Length of the response
	bytes = append(bytes, r.body...)                                  // Data bytes

	return bytes
}
