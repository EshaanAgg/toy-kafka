package response

import "encoding/binary"

type Response struct {
	body []byte
}

func NewResponse(correlationID int32) *Response {
	r := &Response{}
	r.WriteInt32(correlationID)

	return r
}

func (r *Response) Bytes() []byte {
	// Add correlation ID to the response
	l := len(r.body)

	var bytes []byte
	// Write the length of the response
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(l))

	bytes = append(bytes, r.body...)

	return bytes
}
