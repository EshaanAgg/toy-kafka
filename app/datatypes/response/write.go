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
		bytes := make([]byte, 0)

		for {
			toWrite := byte(val & 0x7F)
			val >>= 7
			if val != 0 {
				toWrite |= 0x80
				bytes = append(bytes, toWrite)
			} else {
				bytes = append(bytes, toWrite)
				break
			}
		}

		r.body = append(r.body, bytes...)
	}
}
