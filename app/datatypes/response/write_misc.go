package response

func (r *Response) WriteEmptyTaggedFields() {
	// Tagged fields are represented by a compact array.
	// We set the same to be nil to indicate that there are no tagged fields.
	r.WriteVarUInt(0)
}

func (r *Response) WriteUUID(val []byte) {
	if len(val) != 16 {
		panic("UUID must be 16 bytes")
	}
	r.body = append(r.body, val...)
}
