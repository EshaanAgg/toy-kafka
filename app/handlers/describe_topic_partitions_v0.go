package handlers

import (
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
	// TODO: Add appropriate types for the fields
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

	for _, topic := range r.Body.Topics.Values {
		body.Topics.Append(getDescribeTopicPartitionsV0ResponseTopic(topic))
	}

	header := &protocol.HeaderV1{
		CorrelationID: r.CorrelationID,
	}

	return protocol.NewResponse(header, body).Bytes()
}

func getDescribeTopicPartitionsV0ResponseTopic(topic DescribeTopicPartitionsV0Request_Topic) *DescribeTopicPartitionsV0Response_Topic {
	return &DescribeTopicPartitionsV0Response_Topic{
		ErrorCode: errorcodes.UNKNOWN_TOPIC_OR_PARTITION,
		TopicName: topic.Name,
	}
}
