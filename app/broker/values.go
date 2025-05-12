package broker

import (
	"errors"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

type ValueHeader struct {
	FrameVersion datatypes.Int8
	Type         datatypes.Int8
	Version      datatypes.Int8
}

func newValueHeader(valueBytes []byte) (*ValueHeader, *datatypes.Parser, error) {
	parser := datatypes.NewParser(valueBytes)
	var header ValueHeader
	if err := datatypes.Unmarshal(&header, parser); err != nil {
		return nil, nil, err
	}
	return &header, parser, nil
}

type FeatureLevelValue struct {
	Name         datatypes.CompactString
	Level        datatypes.Int16
	TaggedFields datatypes.TaggedFields
}

type TopicValue struct {
	Name         datatypes.CompactString
	ID           datatypes.UUID
	TaggedFields datatypes.TaggedFields
}

type PartitionValue struct {
	PartitionID      datatypes.Int32
	TopicID          datatypes.UUID
	Replicas         datatypes.CompactArray[datatypes.Int32]
	ISR              datatypes.CompactArray[datatypes.Int32]
	RemovingReplicas datatypes.CompactArray[datatypes.Int32]
	AddingReplicas   datatypes.CompactArray[datatypes.Int32]
	LeaderID         datatypes.Int32
	LeaderEpoch      datatypes.Int32
	PartitionEpoch   datatypes.Int32
	Directories      datatypes.CompactArray[datatypes.UUID]
	TaggedFields     datatypes.TaggedFields
}

func newValueBody[T any](parser *datatypes.Parser) (*T, error) {
	var value T
	if err := datatypes.Unmarshal(&value, parser); err != nil {
		return nil, err
	}

	// Check if the parser is at the end
	if !parser.IsAtEnd() {
		return nil, errors.New("expected the parser to be at the end, but it is not")
	}

	return &value, nil
}
