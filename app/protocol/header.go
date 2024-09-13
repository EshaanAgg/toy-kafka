package protocol

type RequestHeader struct {
	APIKey        int16
	APIVersion    int16
	CorrelationID int32
	ClientID      *string
}

func ParseRequestHeader(r *Request) *RequestHeader {
	apiKey := r.ReadInt16()
	apiVersion := r.ReadInt16()
	correlationID := r.ReadInt32()

	return &RequestHeader{
		APIKey:        apiKey,
		APIVersion:    apiVersion,
		CorrelationID: correlationID,
		ClientID:      nil,
	}
}

func (h *RequestHeader) ValidateAPIVersion() bool {
	return h.APIVersion >= 0 && h.APIVersion <= 4
}
