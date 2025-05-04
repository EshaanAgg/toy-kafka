package response

func (r *Response) WriteCompactArrayLength(length int) {
	// Compact array length is represented by a varint of length + 1.
	r.WriteVarUInt(uint(length + 1))
}

func WriteCompactArray[T any](r *Response, vals []T, writeFunc func(*Response, T)) {
	r.WriteCompactArrayLength(len(vals))
	for _, val := range vals {
		writeFunc(r, val)
	}
}
