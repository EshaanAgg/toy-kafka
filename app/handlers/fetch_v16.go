package handlers

import (
	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
	"github.com/EshaanAgg/toy-kafka/app/datatypes/response"
)

type FetchV16Body struct {
}

type FetchV16Request struct {
	*request.RequestHeader
	Body *FetchV16Body
}

// Fetch Request (Version: 16) => max_wait_ms min_bytes max_bytes isolation_level session_id session_epoch [topics] [forgotten_topics_data] rack_id _tagged_fields
//	max_wait_ms => INT32
//	min_bytes => INT32
//	max_bytes => INT32
//	isolation_level => INT8
//	session_id => INT32
//	session_epoch => INT32
//	topics => topic_id [partitions] _tagged_fields
//	  topic_id => UUID
//	  partitions => partition current_leader_epoch fetch_offset last_fetched_epoch log_start_offset partition_max_bytes _tagged_fields
//	    partition => INT32
//	    current_leader_epoch => INT32
//	    fetch_offset => INT64
//	    last_fetched_epoch => INT32
//	    log_start_offset => INT64
//	    partition_max_bytes => INT32
//	forgotten_topics_data => topic_id [partitions] _tagged_fields
//	  topic_id => UUID
//	  partitions => INT32
//	rack_id => COMPACT_STRING

func NewFetchV16Request(r *request.RequestHeader) (request.Request, error) {
	return &FetchV16Request{
		RequestHeader: r,
		Body:          &FetchV16Body{},
	}, nil
}

func (r *FetchV16Request) Handle() (*response.Response, error) {
	return nil, nil
}
