package handlers

import (
	"fmt"
	"slices"

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

	name, ok := broker.TopicNameFromID[topic.TopicID]
	if !ok {
		// Topic not found, so create a parition with UNKNOWN_TOPIC_ID
		partition := &FetchV16Response_Partition{
			Index:     0,
			ErrorCode: UNKNOWN_TOPIC_ID_ERROR_CODE,
		}
		response.Partitions.Append(partition)
		return response
	}

	// If the topic is found, create partitions for it
	onDiskPartitions, ok := broker.TopicPartitions[name]
	if !ok {
		panic(fmt.Sprintf("Topic %s exists, but there are no partitions for the same", name))
	}

	for _, requestedPartition := range topic.Partitions.Values {
		partition := &FetchV16Response_Partition{
			Index:     requestedPartition.Partition,
			ErrorCode: NO_ERROR_CODE,
		}
		partitionIdx := int(requestedPartition.Partition)
		existsOnDisk := slices.Contains(onDiskPartitions, partitionIdx)

		if !existsOnDisk {
			// Update the error code to indicate that the partition does not exist
			partition.ErrorCode = UNKNOWN_TOPIC_OR_PARTITION_ERROR_CODE
		} else {
			// Fetch the partition data from the broker
			partitionData, err := broker.GetTopicPartitionData(name, partitionIdx)
			if err != nil {
				panic(fmt.Sprintf("Unable to get partition data for topic %s and partition %d: %v", name, partitionIdx, err))
			}
			partition.Records = partitionData
		}

		// Append the partition to the response
		response.Partitions.Append(partition)
	}

	return response
}
