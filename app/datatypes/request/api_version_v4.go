package request

import (
	"fmt"
)

type APIVersionV4Body struct {
	ClientSoftwareName    string
	ClientSoftwareVersion string
}

type APIVersionV4Request struct {
	*RequestHeader
	Body *APIVersionV4Body
}

// ApiVersions Request (Version 4) => client_software_name client_software_version _tagged_fields
//
//	client_software_name => COMPACT_STRING
//	client_software_version => COMPACT_STRING

func NewAPIVersionV4Request(r *RequestHeader) (Request, error) {
	clientSoftwareName, err := r.ReadCompactString()
	if err != nil {
		return nil, fmt.Errorf("newAPIVersionBody.clientSoftwareName: %w", err)
	}

	clientSoftwareVersion, err := r.ReadCompactString()
	if err != nil {
		return nil, fmt.Errorf("newAPIVersionBody.clientSoftwareVersion: %w", err)
	}

	return &APIVersionV4Request{
		RequestHeader: r,
		Body: &APIVersionV4Body{
			ClientSoftwareName:    clientSoftwareName,
			ClientSoftwareVersion: clientSoftwareVersion,
		},
	}, nil
}
