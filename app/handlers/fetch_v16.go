package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

type FetchV16_Partition struct {
	Partition          datatypes.Int32
	CurrentLeaderEpoch datatypes.Int32
	FetchOffset        datatypes.Int64
	LastFetchedEpoch   datatypes.Int32
	LogStartOffset     datatypes.Int64
	PartitionMaxBytes  datatypes.Int32
}

type FetchV16_Topic struct {
	TopicID      datatypes.UUID
	Paritions    datatypes.CompactArray[FetchV16_Partition]
	TaggedFields datatypes.TaggedFields
}

type FetchV16_ForgottenTopic struct {
	TopicID      datatypes.UUID
	Paritions    datatypes.CompactArray[datatypes.Int32]
	TaggedFields datatypes.TaggedFields
}

type FetchV16Body struct {
	MaxWaitMS       datatypes.Int32
	MinBytes        datatypes.Int32
	MaxBytes        datatypes.Int32
	IsolationLevel  datatypes.Int8
	SessionID       datatypes.Int32
	SessionEpoch    datatypes.Int32
	Topics          datatypes.CompactArray[FetchV16_Topic]
	ForgottenTopics datatypes.CompactArray[FetchV16_ForgottenTopic]
	RackID          datatypes.CompactString
	TaggedFields    datatypes.TaggedFields
}

type FetchV16Request struct {
	*protocol.RequestHeader
	Body *FetchV16Body
}

func NewFetch16Request(r *protocol.RequestHeader) (protocol.Request, error) {
	var body FetchV16Body
	if err := datatypes.Unmarshal(&body, r.P); err != nil {
		return nil, fmt.Errorf("unable to decode the request body: %w", err)
	}

	return &FetchV16Request{
		RequestHeader: r,
		Body:          &body,
	}, nil
}
