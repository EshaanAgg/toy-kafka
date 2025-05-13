package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

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
		NewFn:      NewFetchV16Request,
	},
	75: {
		MinVersion: 0,
		MaxVersion: 0,
		NewFn:      NewDescribeTopicPartitionsV0Request,
	},
}
