package protocol

import "github.com/EshaanAgg/toy-kafka/app/datatypes"

type HeaderV0 struct {
	CorrelationID datatypes.Int32
}

type HeaderV1 struct {
	CorrelationID datatypes.Int32
	TaggedFileds  datatypes.TaggedFields
}
