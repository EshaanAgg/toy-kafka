package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

const NO_ERROR_CODE = 0
const UNKNOWN_TOPIC_OR_PARTITION_ERROR_CODE = 3
const UNSUPPORTED_API_VERSION_ERROR_CODE = 35
const UNKNOWN_TOPIC_ID_ERROR_CODE = 100

type SupportedAPI struct {
	MinVersion int16
	MaxVersion int16
	NewFn      func(*protocol.RequestHeader) (protocol.Request, error)
}

var RequestKeyMap = map[datatypes.Int16]SupportedAPI{
	18: {
		MinVersion: 4,
		MaxVersion: 4,
		NewFn:      NewAPIVersionsV4Request,
	},
	1: {
		MinVersion: 16,
		MaxVersion: 16,
		NewFn:      NewFetch16Request,
	},
}
