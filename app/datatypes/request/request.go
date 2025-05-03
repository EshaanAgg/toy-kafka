package request

import "github.com/EshaanAgg/toy-kafka/app/datatypes/response"

// A request is an interface that can be parser from a byte array.
// The various request types implement this interface.
type Request interface {
	Handle() (*response.Response, error)
}
