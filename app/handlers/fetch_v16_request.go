package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
)

// Request
type FetchV16Body struct {
	MaxWaitMS       int32                 `kafka:"int32"`
	MinBytes        int32                 `kafka:"int32"`
	MaxBytes        int32                 `kafka:"int32"`
	IsolationLevel  int8                  `kafka:"int8"`
	SessionID       int32                 `kafka:"int32"`
	SessionEpoch    int32                 `kafka:"int32"`
	Topics          []*FetchV16Body_Topic `kafka:"compact_array:struct"`
	ForgottenTopics []*FetchV16Body_Topic `kafka:"compact_array:struct"`
	RackID          string                `kafka:"compact_string"`
}

type FetchV16Body_Topic struct {
	TopicID    string                    `kafka:"uuid"`
	Partitions []*FetchV16Body_Partition `kafka:"compact_array:struct"`
}

type FetchV16Body_Partition struct {
	Partition          int32 `kafka:"int32"`
	CurrentLeaderEpoch int32 `kafka:"int32"`
	FetchOffset        int64 `kafka:"int64"`
	LastFetchedEpoch   int32 `kafka:"int32"`
	LogStartOffset     int64 `kafka:"int64"`
	PartitionMaxBytes  int32 `kafka:"int32"`
}

type FetchV16Request struct {
	*request.RequestHeader
	Body *FetchV16Body
}

func NewFetchV16Request(r *request.RequestHeader) (request.Request, error) {
	return &FetchV16Request{
		RequestHeader: r,
		Body:          &FetchV16Body{},
	}, nil
}
