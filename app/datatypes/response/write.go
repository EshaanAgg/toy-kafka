package response

import "encoding/binary"

func (r *Response) WriteInt16(vals ...int16) {
	for _, val := range vals {
		r.body = binary.BigEndian.AppendUint16(r.body, uint16(val))
	}
}

func (r *Response) WriteInt32(vals ...int32) {
	for _, val := range vals {
		r.body = binary.BigEndian.AppendUint32(r.body, uint32(val))
	}
}

func (r *Response) WriteVarInt(vals ...int) {
	for _, val := range vals {
		r.body = binary.AppendVarint(r.body, int64(val))
	}
}

func (r *Response) WriteVarUInt(vals ...uint) {
	for _, val := range vals {
		r.body = binary.AppendUvarint(r.body, uint64(val))
	}
}

func (r *Response) WriteEmptyTaggedFields() {
	// Tagged fields are represented by a compact array.
	// We set the same to be nil to indicate that there are no tagged fields.
	r.WriteVarUInt(0)
}

func (r *Response) WriteCompactArrayLength(length int) {
	// Compact array length is represented by a varint of length + 1.
	r.WriteVarInt(length + 1)
}

func WriteCompactArray[T any](r *Response, vals []T, writeFunc func(*Response, T)) {
	r.WriteCompactArrayLength(len(vals))
	for _, val := range vals {
		writeFunc(r, val)
	}
}
