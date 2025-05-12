package broker

const CLUSTER_METADATA_FILE = "/tmp/kraft-combined-logs/__cluster_metadata-0/00000000000000000000.log"

type Broker struct {
}

func (b *Broker) GetClusterMetadata() {
	// Read the cluster metadata file
	_, err := readLogFile(CLUSTER_METADATA_FILE)
	if err != nil {
		panic(err)
	}

}
