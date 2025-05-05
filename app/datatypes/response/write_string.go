package response

func (r *Response) WriteString(val string) {
	r.WriteInt16(int16(len(val)))   // Length of the string as INT16
	r.body = append(r.body, val...) // Data bytes of the string
}

func (r *Response) WriteNullableString(val *string) {
	if val == nil {
		r.WriteInt16(-1) // -1 for null value
	} else {
		r.WriteInt16(int16(len(*val)))   // Length of the string
		r.body = append(r.body, *val...) // Data bytes of the string
	}
}

func (r *Response) WriteCompactString(val string) {
	r.WriteVarUInt(uint(len(val) + 1)) // Length of the string + 1 as VARUINT
	r.body = append(r.body, val...)    // Data bytes of the string
}
