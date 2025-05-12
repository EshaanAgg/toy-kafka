package broker

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

type RecordBatch struct {
	BaseOffset           datatypes.Int64
	BatchLength          datatypes.Int32
	PartitionLeaderEpoch datatypes.Int32
	MagicByte            datatypes.Int8
	CRC                  datatypes.UInt32
	Attributes           datatypes.Int16
	LastOffsetDelta      datatypes.Int32
	BaseTimestamp        datatypes.Int64
	MaxTimestamp         datatypes.Int64
	ProducerID           datatypes.Int64
	ProducerEpoch        datatypes.Int16
	BaseSequence         datatypes.Int32
	Records              datatypes.Array[Record]
}

func (b *RecordBatch) assertNoCompression() {
	// Attributes is 16 bits, we need to access bits 0-2
	compressionMask := datatypes.Int16(0b0000000000000111)
	compression := b.Attributes & compressionMask
	if compression != 0 {
		panic(fmt.Sprintf("reading compressed log files is not supported, recieved compression %d", compression))
	}
}

type Record struct {
	Length         datatypes.VarInt
	Attributes     datatypes.Int8
	TimestampDelta datatypes.VarLong
	OffsetDelta    datatypes.VarInt
	Key            datatypes.VarIntBytes
	Value          datatypes.VarIntBytes
	Headers        datatypes.VarIntArray[Header]
}

type Header struct {
	Key   datatypes.VarIntString
	Value datatypes.VarIntBytes
}
