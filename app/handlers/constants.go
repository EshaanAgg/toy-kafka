package handlers

import "github.com/EshaanAgg/toy-kafka/app/datatypes/request"

const UNSUPPORTED_API_VERSION_ERROR_CODE = 35

type SupportedAPI struct {
	MinVersion int16
	MaxVersion int16
	NewFn      func(*request.RequestHeader) (request.Request, error)
}

var RequestKeyMap = map[int16]SupportedAPI{
	1: {
		MinVersion: 16,
		MaxVersion: 16,
		NewFn:      NewFetchV16Request,
	},
	18: {
		MinVersion: 4,
		MaxVersion: 4,
		NewFn:      NewAPIVersionsV4Request,
	},
}
