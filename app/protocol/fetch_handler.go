package protocol

type FetchDataRequest struct {
	ReplicaID           int32
	MaxWaitMS           int32
	MinBytes            int32
	MaxBytes            int32
	IsolationLevel      int8
	SessionID           int32
	SessionEpoch        int32
	Topics              []FetchTopic
	ForgottenTopicsData []ForgottenTopicData
	RackID              string
}

func NewFetchRequest(r *Request) *FetchDataRequest {
	f := &FetchDataRequest{}
	f.Decode(r)
	return f
}

func (f *FetchDataRequest) Decode(r *Request) {
	f.ReplicaID = r.ReadInt32()
	f.MaxWaitMS = r.ReadInt32()
	f.MinBytes = r.ReadInt32()
	f.MaxBytes = r.ReadInt32()
	f.IsolationLevel = r.ReadInt8()
	f.SessionID = r.ReadInt32()
	f.SessionEpoch = r.ReadInt32()

	topicCnt := r.ReadVarInt() - 1
	f.Topics = make([]FetchTopic, topicCnt)
	for i := 0; i < topicCnt; i++ {
		f.Topics[i].Decode(r)
	}

	forgottenTopicCnt := r.ReadVarInt() - 1
	f.ForgottenTopicsData = make([]ForgottenTopicData, forgottenTopicCnt)
	for i := 0; i < forgottenTopicCnt; i++ {
		f.ForgottenTopicsData[i].Decode(r)
	}

	f.RackID = r.ReadCompactString()

	r.ReadTaggedFields()
}

type FetchTopic struct {
	TopicID    string
	Partitions []FetchPartition
}

func (t *FetchTopic) Decode(r *Request) {
	t.TopicID = r.ReadUUID()

	partitionCnt := r.ReadVarInt() - 1
	t.Partitions = make([]FetchPartition, partitionCnt)
	for i := 0; i < partitionCnt; i++ {
		t.Partitions[i].Decode(r)
	}

	r.ReadTaggedFields()
}

type FetchPartition struct {
	PartitionID        int32
	CurrentLeaderEpoch int32
	FetchOffset        int64
	LastFetchedEpoch   int32
	LogStartOffset     int64
	PartitionMaxBytes  int32
}

func (p *FetchPartition) Decode(r *Request) {
	p.PartitionID = r.ReadInt32()
	p.CurrentLeaderEpoch = r.ReadInt32()
	p.FetchOffset = r.ReadInt64()
	p.LastFetchedEpoch = r.ReadInt32()
	p.LogStartOffset = r.ReadInt64()
	p.PartitionMaxBytes = r.ReadInt32()

	r.ReadTaggedFields()
}

type ForgottenTopicData struct {
	TopicID    string
	Partitions []int32
}

func (f *ForgottenTopicData) Decode(r *Request) {
	f.TopicID = r.ReadUUID()
	f.Partitions = DecodeInt32Array(r)

	r.ReadTaggedFields()
}

type FetchDataBody struct {
	ThrottleTimeMS int32
	ErrorCode      int16
	SessionID      int32
	Responses      []FetchDataBodyResponse
}

func NewFetchDataBody(r *FetchDataRequest) *FetchDataBody {
	return &FetchDataBody{
		ThrottleTimeMS: 0,
		ErrorCode:      0,
		SessionID:      r.SessionID,
		// TODO: Implement getting responses from the request
		Responses: []FetchDataBodyResponse{},
	}
}

func (b *FetchDataBody) Encode(r *Response) {
	r.WriteInt32(b.ThrottleTimeMS)
	r.WriteInt16(b.ErrorCode)
	r.WriteInt32(b.SessionID)

	r.WriteVarInt(len(b.Responses) + 1)
	for _, resp := range b.Responses {
		resp.Encode(r)
	}
}

// TODO: Implement response parsing and encoding
type FetchDataBodyResponse struct {
	Topic string
}

func (r *FetchDataBodyResponse) Encode(res *Response) {}

func GetFetchBody(r *Request) *FetchDataBody {
	req := NewFetchRequest(r)
	return NewFetchDataBody(req)
}
