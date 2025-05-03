package request

import "github.com/EshaanAgg/toy-kafka/app/datatypes/response"

// ApiVersions Response (Version: 4) => error_code [api_keys] throttle_time_ms _tagged_fields
//
//	error_code => INT16
//	api_keys => api_key min_version max_version _tagged_fields
//	  api_key => INT16
//	  min_version => INT16
//	  max_version => INT16
//	throttle_time_ms => INT32

func (r *APIVersionV4Request) Handle() (*response.Response, error) {
	var res response.Response

	res.WriteInt16(0) // Error code
	for key, api := range RequestKeyMap {
		res.WriteInt16(key, api.MinVersion, api.MaxVersion)
	}
	res.WriteInt16(0) // Throttle time

	return &res, nil
}
