package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/broker"
	"github.com/EshaanAgg/toy-kafka/app/datatypes"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
	"github.com/EshaanAgg/toy-kafka/app/handlers/errorcodes"
)

type DescribeTopicPartitions_Cursor struct {
	TopicName      datatypes.CompactString
	PartitionIndex datatypes.Int32
	TaggedFields   datatypes.TaggedFields
}

type DescribeTopicPartitionsV0Request struct {
	*protocol.RequestHeader
	Body *DescribeTopicPartitionsV0Body
}

type DescribeTopicPartitionsV0Request_Topic struct {
	Name         datatypes.CompactString
	TaggedFields datatypes.TaggedFields
}

type DescribeTopicPartitionsV0Body struct {
	Topics                 datatypes.CompactArray[DescribeTopicPartitionsV0Request_Topic]
	ResponsePartitionLimit datatypes.Int32
	Cursor                 datatypes.Nullable[DescribeTopicPartitions_Cursor]
	TaggedFields           datatypes.TaggedFields
}

type DescribeTopicPartitionsV0Response_Partition struct {
	ErrorCode              datatypes.Int16
	PartitionIndex         datatypes.Int32
	LeaderID               datatypes.Int32
	LeaderEpoch            datatypes.Int32
	ReplicaNodes           datatypes.CompactArray[datatypes.Int32]
	ISRNodes               datatypes.CompactArray[datatypes.Int32]
	EligibleLeaderReplicas datatypes.CompactArray[datatypes.Int32]
	LastKnownELR           datatypes.CompactArray[datatypes.Int32]
	OfflineReplicas        datatypes.CompactArray[datatypes.Int32]
	TaggedFields           datatypes.TaggedFields
}

type DescribeTopicPartitionsV0Response_Topic struct {
	ErrorCode                 datatypes.Int16
	TopicName                 datatypes.CompactString
	TopicID                   datatypes.UUID
	IsInternal                datatypes.Boolean
	Partitions                datatypes.CompactArray[DescribeTopicPartitionsV0Response_Partition]
	TopicAuthorizedOperations datatypes.Int32
	TaggedFields              datatypes.TaggedFields
}

type DescribeTopicPartitionsV0Response struct {
	ThrottleTimeMs datatypes.Int32
	Topics         datatypes.CompactArray[DescribeTopicPartitionsV0Response_Topic]
	NextCursor     datatypes.Nullable[DescribeTopicPartitions_Cursor]
	TaggedFields   datatypes.TaggedFields
}

func NewDescribeTopicPartitionsV0Request(r *protocol.RequestHeader) (protocol.Request, error) {
	body, err := getBody[DescribeTopicPartitionsV0Body](r.P)
	if err != nil {
		return nil, err
	}

	return &DescribeTopicPartitionsV0Request{
		RequestHeader: r,
		Body:          body,
	}, nil
}

func (r *DescribeTopicPartitionsV0Request) Handle() ([]byte, error) {
	body := &DescribeTopicPartitionsV0Response{
		ThrottleTimeMs: 0,
		NextCursor:     datatypes.NewNullable[DescribeTopicPartitions_Cursor](),
	}

	broker, err := broker.NewBroker()
	if err != nil {
		return nil, fmt.Errorf("unable to create broker: %w", err)
	}

	for _, topic := range r.Body.Topics.Values {
		body.Topics.Append(getDescribeTopicPartitionsV0ResponseTopic(topic, broker))
	}

	header := &protocol.HeaderV1{
		CorrelationID: r.CorrelationID,
	}

	return protocol.NewResponse(header, body).Bytes()
}

func getDescribeTopicPartitionsV0ResponseTopic(
	topic DescribeTopicPartitionsV0Request_Topic,
	broker *broker.Broker,
) *DescribeTopicPartitionsV0Response_Topic {
	response := &DescribeTopicPartitionsV0Response_Topic{
		TopicName: topic.Name,
		ErrorCode: errorcodes.NO_ERROR,
	}

	topicID, ok := broker.TopicIDFromName[string(topic.Name)]
	if !ok {
		response.ErrorCode = errorcodes.UNKNOWN_TOPIC_OR_PARTITION
		return response
	}

	response.TopicID = topicID

	partitionIndexes := broker.TopicPartitions[string(topic.Name)]
	for _, partitionIndex := range partitionIndexes {
		p := &DescribeTopicPartitionsV0Response_Partition{
			PartitionIndex: datatypes.Int32(partitionIndex),
			ErrorCode:      errorcodes.NO_ERROR,
		}
		response.Partitions.Append(p)
	}

	return response
}
