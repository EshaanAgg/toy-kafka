package handlers

import (
	"fmt"

	"github.com/EshaanAgg/toy-kafka/app/datatypes"
)

// Returns the body of the request as a pointer to the type T.
// It also checks if there are any extra bytes in the request body after
// unmarshalling the body. If there are, it returns an error.
func getBody[T any](p *datatypes.Parser) (*T, error) {
	var body T
	if err := datatypes.Unmarshal(&body, p); err != nil {
		return nil, fmt.Errorf("unable to decode the request body: %w", err)
	}
	if !p.IsAtEnd() {
		p.Debug()
		return nil, fmt.Errorf("extra bytes in request body")
	}
	return &body, nil
}
