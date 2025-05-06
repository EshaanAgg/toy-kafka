package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

type APIVersionsV4Request struct {
	*protocol.RequestHeader
	Body *APIVersionsV4Body
}
type APIVersionsV4Body struct {
	ClientSoftwareName    datatypes.CompactString
	ClientSoftwareVersion datatypes.CompactString
}

func NewAPIVersionsV4Request(r *protocol.RequestHeader) (protocol.Request, error) {
	var body APIVersionsV4Body
	if err := datatypes.Unmarshal(&body, r.P); err != nil {
		return nil, fmt.Errorf("unable to decode the request body: %w", err)
	}

	return &APIVersionsV4Request{
		RequestHeader: r,
		Body:          &body,
	}, nil
}

func (r *APIVersionsV4Request) getErrorCode() datatypes.Int16 {
	api := RequestKeyMap[r.APIKey]
	if api.MinVersion > int16(r.APIVersion) || api.MaxVersion < int16(r.APIVersion) {
		return UNSUPPORTED_API_VERSION_ERROR_CODE
	}
	return 0
}

type APIVersionsV4Response_APIKey struct {
	APIKey       datatypes.Int16
	MinVersion   datatypes.Int16
	MaxVersion   datatypes.Int16
	TaggedFields datatypes.TaggedFields
}

type APIVersionsV4Response struct {
	ErrorCode    datatypes.Int16
	APIKeys      datatypes.CompactArray[APIVersionsV4Response_APIKey]
	ThrottleTime datatypes.Int32
	TaggedFields datatypes.TaggedFields
}

func (r *APIVersionsV4Request) Handle() ([]byte, error) {
	body := &APIVersionsV4Response{
		ErrorCode:    r.getErrorCode(),
		ThrottleTime: 0,
	}

	for key, api := range RequestKeyMap {
		body.APIKeys.Values = append(body.APIKeys.Values, APIVersionsV4Response_APIKey{
			APIKey:     datatypes.Int16(key),
			MinVersion: datatypes.Int16(api.MinVersion),
			MaxVersion: datatypes.Int16(api.MaxVersion),
		})
	}

	header := &protocol.HeaderV0{
		CorrelationID: r.CorrelationID,
	}

	return protocol.NewResponse(header, body).Bytes()
}
