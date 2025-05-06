package datatypes

import (
	"bytes"
	"fmt"
	"reflect"
)

// Marshal takes in a struct and returns a byte slice obtained due to
// marshalling the struct into bytes. All the fields of the struct must be
// of type MarshableDataType, and no nesting is allowed.
func Marshal(obj any, b *bytes.Buffer) error {
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %T", obj)
	}

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	for i := range v.NumField() {
		field := v.Field(i)

		// Expect the field to implement the MarshableDataType interface
		if !field.CanInterface() {
			return fmt.Errorf("field %s is not exportable", t.Field(i).Name)
		}

		dt, ok := field.Interface().(MarshableDataType)
		if !ok {
			return fmt.Errorf("expected MarshableDataType, got %T", field.Interface())
		}

		if err := dt.Marshal(b); err != nil {
			return fmt.Errorf("%s: %w", t.Field(i).Name, err)
		}
	}

	return nil
}

// Unmarshal must be called with a pointer to a struct where
// the unmarshalled data is to be stored.
func Unmarshal(obj any, p *Parser) error {
	// We need a pointer to struct for unmarshalling
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to struct, got %T", obj)
	}

	// Dereference the pointer to get the struct value
	structVal := objVal.Elem()
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got pointer to %s", structVal.Kind())
	}

	structType := structVal.Type()

	for i := range structVal.NumField() {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// Check if the field is exportable (public)
		if !field.CanSet() {
			return fmt.Errorf("field %s is not exportable", fieldType.Name)
		}

		fieldInterface := field.Addr().Interface()
		dt, ok := fieldInterface.(UnmarshableDataType)
		if !ok {
			return fmt.Errorf("field %s (%T) does not implement UnmarshableDataType", fieldType.Name, field.Interface())
		}

		if err := dt.Unmarshal(p); err != nil {
			return fmt.Errorf("%s: %w", fieldType.Name, err)
		}
	}

	return nil
}
