package handlers

import "github.com/EshaanAgg/toy-kafka/app/datatypes/response"

type FetchV16Response struct {
	ThrottleTimeMS int32                       `kafka:"int32"`
	ErrorCode      int16                       `kafka:"int16"`
	SessionID      int32                       `kafka:"int32"`
	Responses      []FetchV16Response_Response `kafka:"compact_array:struct"`
}

type FetchV16Response_Response struct {
	TopicID    string                       `kafka:"uuid"`
	Partitions []FetchV16Response_Partition `kafka:"compact_array:struct"`
}

type FetchV16Response_Partition struct {
	PartitionIndex       int32                                 `kafka:"int32"`
	ErrorCode            int16                                 `kafka:"int16"`
	HighWatermark        int64                                 `kafka:"int64"`
	LastStableOffset     int64                                 `kafka:"int64"`
	LogStartOffset       int64                                 `kafka:"int64"`
	AbortedTransactions  []FetchV16Response_AbortedTransaction `kafka:"compact_array:struct"`
	PreferredReadReplica int32                                 `kafka:"int32"`
	Records              []byte                                `kafka:"compact_records"`
}

type FetchV16Response_AbortedTransaction struct {
	ProducerID  int64 `kafka:"int64"`
	FirstOffset int64 `kafka:"int64"`
}

func (r *FetchV16Request) Handle() (*response.Response, error) {
	body := &FetchV16Response{}
	return response.NewResponseWithBody(r.CorrelationID, body)
}
