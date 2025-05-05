package response

import (
	"encoding/binary"
)

func (r *Response) WriteInt8(vals ...int8) {
	for _, val := range vals {
		r.body = append(r.body, byte(val))
	}
}

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

func (r *Response) WriteInt64(vals ...int64) {
	for _, val := range vals {
		r.body = binary.BigEndian.AppendUint64(r.body, uint64(val))
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
