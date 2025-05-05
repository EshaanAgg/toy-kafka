package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/response"
)

// Request

type APIVersionsV4Request struct {
	*request.RequestHeader
	Body *APIVersionsV4Body
}
type APIVersionsV4Body struct {
	ClientSoftwareName    string `kafka:"compact_string"`
	ClientSoftwareVersion string `kafka:"compact_string"`
}

func NewAPIVersionsV4Request(r *request.RequestHeader) (request.Request, error) {
	var body APIVersionsV4Body
	if err := r.AutoDecodeBody("APIVersionsV4Request.Body", &body, true); err != nil {
		return nil, err
	}

	return &APIVersionsV4Request{
		RequestHeader: r,
		Body:          &body,
	}, nil
}

// Response

type APIVersionsV4Response_APIKey struct {
	APIKey     int16 `kafka:"int16"`
	MinVersion int16 `kafka:"int16"`
	MaxVersion int16 `kafka:"int16"`
}

type APIVersionsV4Response struct {
	ErrorCode    int16                          `kafka:"int16"`
	APIKeys      []APIVersionsV4Response_APIKey `kafka:"compact_array:struct"`
	ThrottleTime int32                          `kafka:"int32"`
}

func (r *APIVersionsV4Request) Handle() (*response.Response, error) {
	body := &APIVersionsV4Response{
		ErrorCode:    r.getErrorCode(),
		APIKeys:      make([]APIVersionsV4Response_APIKey, 0),
		ThrottleTime: 0,
	}

	for key, api := range RequestKeyMap {
		body.APIKeys = append(body.APIKeys, APIVersionsV4Response_APIKey{
			APIKey:     key,
			MinVersion: api.MinVersion,
			MaxVersion: api.MaxVersion,
		})
	}

	return response.NewResponseWithBody(r.CorrelationID, body)
}

func (r *APIVersionsV4Request) getErrorCode() int16 {
	api := RequestKeyMap[r.APIKey]
	if api.MinVersion > r.APIVersion || api.MaxVersion < r.APIVersion {
		return UNSUPPORTED_API_VERSION_ERROR_CODE
	}
	return 0
}
