package broker

import (
	"fmt"
	"os"
	"strings"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

const LOG_BASE_DIR = "/tmp/kraft-combined-logs"
const CLUSTER_METADATA_FILE = LOG_BASE_DIR + "/__cluster_metadata-0/00000000000000000000.log"

type Broker struct {
	TopicNameFromID map[datatypes.UUID]string
	TopicPartitions map[string][]int
}

func NewBroker() (*Broker, error) {
	b := &Broker{
		TopicNameFromID: make(map[datatypes.UUID]string),
		TopicPartitions: make(map[string][]int),
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
			partitionValue, err := newValueBody[PartitionValue](p)
			if err != nil {
				return fmt.Errorf("unable to decode the partition value: %w", err)
			}
			b.processPartitionValue(partitionValue)

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

func (b *Broker) processPartitionValue(partitionValue *PartitionValue) {
	// Store the partition ID in the map using the topic name as the key
	topicName, ok := b.TopicNameFromID[partitionValue.TopicID]
	if !ok {
		fmt.Printf("Topic ID %s not found in TopicNameFromID map\n", partitionValue.TopicID)
		return
	}
	b.TopicPartitions[topicName] = append(b.TopicPartitions[topicName], int(partitionValue.PartitionID))
}

// GetTopicPartitionData returns the data for a given topic and partition ID.
// It reads all the log files in the topic's partition directory and concatenates their contents.
func (b *Broker) GetTopicPartitionData(topicName string, partitionID int) ([]byte, error) {
	folderPath := fmt.Sprintf("%s/%s-%d", LOG_BASE_DIR, topicName, partitionID)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder '%s' does not exist", folderPath)
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory '%s': %w", folderPath, err)
	}
	if len(files) == 0 {
		return []byte{}, nil
	}

	d := make([]byte, 0)
	for _, file := range files {
		// Only read from .log files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", folderPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("unable to read file '%s': %w", filePath, err)
		}
		d = append(d, data...)
	}

	return d, nil
}
