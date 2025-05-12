package broker

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

const CLUSTER_METADATA_FILE = "/tmp/kraft-combined-logs/__cluster_metadata-0/00000000000000000000.log"

type Broker struct {
	TopicNameFromID map[datatypes.UUID]string
}

func NewBroker() (*Broker, error) {
	b := &Broker{
		TopicNameFromID: make(map[datatypes.UUID]string),
	}
	// If the cluster metadata file does not exist, return an empty broker
	if _, err := os.Stat(CLUSTER_METADATA_FILE); os.IsNotExist(err) {
		return b, nil
	}

	// Read the cluster metadata file and process all the batches
	batches, err := readLogFile(CLUSTER_METADATA_FILE)
	if err != nil {
		return nil, fmt.Errorf("unable to read cluster metadata file: %w", err)
	}
	for _, batch := range batches {
		if err := b.processBatch(batch); err != nil {
			return nil, fmt.Errorf("unable to process batch: %w", err)
		}
	}

	return b, nil
}

func (b *Broker) processBatch(batch *RecordBatch) error {
	// Process the batch and extract topic names
	for _, record := range batch.Records.Values {
		header, p, err := newValueHeader(record.Value)
		if err != nil {
			return fmt.Errorf("unable to decode the value header: %w", err)
		}

		switch header.Type {
		case 2:
			// Topic
			topicValue, err := newValueBody[TopicValue](p)
			if err != nil {
				return fmt.Errorf("unable to decode the topic value: %w", err)
			}
			b.processTopicValue(topicValue)

		case 3:
			// Partition
			_, err := newValueBody[PartitionValue](p)
			if err != nil {
				return fmt.Errorf("unable to decode the partition value: %w", err)
			}

		case 12:
			// Feature Level
			_, err := newValueBody[FeatureLevelValue](p)
			if err != nil {
				return fmt.Errorf("unable to decode the feature level value: %w", err)
			}

		default:
			return fmt.Errorf("unknown value type: %d", header.Type)
		}
	}

	return nil
}

func (b *Broker) processTopicValue(topicValue *TopicValue) {
	// Store the topic name in the map using the UUID as the key
	b.TopicNameFromID[topicValue.ID] = string(topicValue.Name)
}
