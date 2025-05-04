package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes/response"
)

const UNSUPPORTED_API_VERSION_ERROR_CODE = 35

// ApiVersions Response (Version: 4) => error_code [api_keys] throttle_time_ms _tagged_fields
//
//	error_code => INT16
//	api_keys => api_key min_version max_version _tagged_fields
//	  api_key => INT16
//	  min_version => INT16
//	  max_version => INT16
//	throttle_time_ms => INT32

func (r *APIVersionsV4Request) Handle() (*response.Response, error) {
	res := response.NewResponse(r.CorrelationID)

	res.WriteInt16(r.getErrorCode()) // Error code

	// API keys
	res.WriteCompactArrayLength(len(RequestKeyMap))
	for key, api := range RequestKeyMap {
		res.WriteInt16(key, api.MinVersion, api.MaxVersion)
		res.WriteEmptyTaggedFields()
	}

	res.WriteInt32(0) // Throttle time
	res.WriteEmptyTaggedFields()

	return res, nil
}

func (r *APIVersionsV4Request) getErrorCode() int16 {
	api := RequestKeyMap[r.APIKey]
	if api.MinVersion > r.APIVersion || api.MaxVersion < r.APIVersion {
		return UNSUPPORTED_API_VERSION_ERROR_CODE
	}
	return 0
}
