package broker

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

// Returns the log file as a an array of RecordBatch.
func readLogFile(filePath string) ([]*RecordBatch, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	parser := datatypes.NewParser(fileBytes)

	batches := make([]*RecordBatch, 0)
	for !parser.IsAtEnd() {
		var batch RecordBatch
		if err := datatypes.Unmarshal(&batch, parser); err != nil {
			return nil, fmt.Errorf("unable to decode the record batch: %w", err)
		}
		batch.assertNoCompression()
		batches = append(batches, &batch)
	}

	return batches, nil
}
