package response

func (r *Response) WriteEmptyTaggedFields() {
	// Tagged fields are represented by a compact array.
	// We set the same to be nil to indicate that there are no tagged fields.
	r.WriteVarUInt(0)
}
