package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/broker"
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
	header := &protocol.HeaderV1{
		CorrelationID: r.CorrelationID,
	}

	resBody := &FetchV16Response{
		ThrottleTimeMS: 0,
		ErrorCode:      NO_ERROR_CODE,
		SessionID:      0,
	}

	broker, err := broker.NewBroker()
	if err != nil {
		return nil, fmt.Errorf("unable to create broker: %w", err)
	}

	for _, topic := range r.Body.Topics.Values {
		resBody.Responses.Append(getResponseForTopic(&topic, broker))
	}

	return protocol.NewResponse(header, resBody).Bytes()
}

// getResponseForTopic creates a FetchV16Response_Response for the given topic.
// If the topic is not found, it creates a partition with UNKNOWN_TOPIC_ID.
func getResponseForTopic(topic *FetchV16_Topic, broker *broker.Broker) *FetchV16Response_Response {
	response := &FetchV16Response_Response{
		TopicID: topic.TopicID,
	}

	_, ok := broker.TopicNameFromID[topic.TopicID]
	if !ok {
		// Topic not found, so create a parition with UNKNOWN_TOPIC_ID
		partition := &FetchV16Response_Partition{
			Index:     0,
			ErrorCode: UNKNOWN_TOPIC_ID_ERROR_CODE,
		}
		response.Partitions.Append(partition)
		return response
	}

	// If the topic is found, create 1 default partition for the same
	partition := &FetchV16Response_Partition{
		Index:     0,
		ErrorCode: NO_ERROR_CODE,
	}
	response.Partitions.Append(partition)

	return response
}
