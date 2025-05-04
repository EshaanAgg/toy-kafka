package request

import (
	"reflect"
)

// AutoDecodeBody decodes the request body based on the tags in the struct fields.
// All fields in the struct must be tagged with "kafka:<tag>". In case of any errors, it returns
// a descriptive error message including the exact field that caused the error.
func (r *RequestHeader) AutoDecodeBody(requestName string, ptr any, isInline bool) error {
	ptrVal := reflect.ValueOf(ptr)

	tag := "struct"
	if isInline {
		tag = "inline_struct"
	}
	return decode(r, tag, ptrVal, requestName)
}
