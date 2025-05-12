package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes/protocol"
)

type DescribeTopicPartitionsV0Request struct {
	*protocol.RequestHeader
	Body *DescribeTopicPartitionsV0Body
}

type DescribeTopicPartitionsV0Body struct{}

type DescribeTopicPartitionsV0Response struct{}

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
	body := &DescribeTopicPartitionsV0Response{}

	header := &protocol.HeaderV1{
		CorrelationID: r.CorrelationID,
	}

	return protocol.NewResponse(header, body).Bytes()
}
