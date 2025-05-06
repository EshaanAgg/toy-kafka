package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

type FetchV16Request_AbortedTransaction struct {
	ProducerID   datatypes.Int64
	FirstOffset  datatypes.Int64
	TaggedFields datatypes.TaggedFields
}

type FetchV16Response_Partition struct {
	Index                datatypes.Int32
	ErrorCode            datatypes.Int16
	HighWatermark        datatypes.Int64
	LastStableOffset     datatypes.Int64
	LogStartOffset       datatypes.Int64
	AbortedTransactions  datatypes.CompactArray[FetchV16Request_AbortedTransaction]
	PreferredReadReplica datatypes.Int32
	Records              datatypes.CompactRecords
	TaggedFields         datatypes.TaggedFields
}

type FetchV16Response_Response struct {
	TopicID      datatypes.UUID
	Partitions   datatypes.CompactArray[FetchV16Response_Partition]
	TaggedFields datatypes.TaggedFields
}

type FetchV16Response struct {
	ThrottleTimeMS datatypes.Int32
	ErrorCode      datatypes.Int16
	SessionID      datatypes.Int32
	Responses      datatypes.CompactArray[FetchV16Response_Response]
	TaggedFields   datatypes.TaggedFields
}

func (r *FetchV16Request) Handle() ([]byte, error) {
	body := &FetchV16Response{}
	header := &protocol.HeaderV1{
		CorrelationID: r.CorrelationID,
	}
	return protocol.NewResponse(header, body).Bytes()
}
