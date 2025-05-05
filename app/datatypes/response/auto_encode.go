package response

import (
	"fmt"
	"reflect"
)

// autoEncodeBody encodes the request body based on the tags in the struct fields.
// All fields in the struct must be tagged with "kafka:<tag>". In case of any errors, it returns
// a descriptive error message including the exact field that caused the error.
func (r *Response) AutoEncodeBody(ptr any, isInline bool) error {
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr || ptrVal.IsNil() {
		return fmt.Errorf("autoEncodeBody: expected non-nil pointer")
	}
	tag := "struct"
	if isInline {
		tag = "inline_struct"
	}
	return encode(r, tag, ptrVal.Elem(), reflect.TypeOf(ptr).Elem().Name())
}
