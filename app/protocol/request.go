package protocol

import (
	"fmt"
	"net"
)

type Request struct {
	bytes           []byte
	byteParseOffset uint

	Conn   *net.Conn
	Length int32
	Header *RequestHeader
}

// Dispatch different handlers for different APIs
const FETCH_API_KEY = 1
const API_VERSIONS_KEY = 18

type ResponseBody interface {
	Encode(r *Response)
}

func (r *Request) Handle() {
	var body ResponseBody

	switch r.Header.APIKey {
	case FETCH_API_KEY:
		body = GetFetchBody(r)

	case API_VERSIONS_KEY:
		body = GetAPIVersionBody(r)

	default:
		fmt.Printf("No handler has been implemented for the request API Key %d", r.Header.APIKey)
		return
	}

	resp := NewResponse()

	resp.WriteInt32(r.Header.CorrelationID) // Header -> CorrelationID
	body.Encode(resp)                       // Body Bytes
	resp.Send(r.Conn)
}
