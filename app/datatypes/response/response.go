package response

import "encoding/binary"

type Response struct {
	body []byte
}

func (r *Response) Bytes() []byte {
	l := len(r.body)

	var bytes []byte
	// Write the length of the response
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(l))
	bytes = append(bytes, r.body...)

	return bytes
}
