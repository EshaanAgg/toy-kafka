package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes/request"
)

type APIVersionsV4Body struct {
	ClientSoftwareName    string
	ClientSoftwareVersion string
}

type APIVersionsV4Request struct {
	*request.RequestHeader
	Body *APIVersionsV4Body
}

// ApiVersions Request (Version 4) => client_software_name client_software_version _tagged_fields
//
//	client_software_name => COMPACT_STRING
//	client_software_version => COMPACT_STRING

func NewAPIVersionsV4Request(r *request.RequestHeader) (request.Request, error) {
	clientSoftwareName, err := r.ReadCompactString()
	if err != nil {
		return nil, fmt.Errorf("newAPIVersionBody.clientSoftwareName: %w", err)
	}

	clientSoftwareVersion, err := r.ReadCompactString()
	if err != nil {
		return nil, fmt.Errorf("newAPIVersionBody.clientSoftwareVersion: %w", err)
	}

	return &APIVersionsV4Request{
		RequestHeader: r,
		Body: &APIVersionsV4Body{
			ClientSoftwareName:    clientSoftwareName,
			ClientSoftwareVersion: clientSoftwareVersion,
		},
	}, nil
}
