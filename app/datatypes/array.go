package datatypes

import (
	"bytes"
	"fmt"
	"reflect"
)

type CompactArray[T any] struct {
	Values []T
}

func (a CompactArray[T]) Marshal(b *bytes.Buffer) error {
	if a.Values == nil {
		var n VarUInt = 0
		if err := n.Marshal(b); err != nil {
			return fmt.Errorf("compact_array.nil_array: %w", err)
		}
		return nil
	}

	l := len(a.Values)

	// Write n + 1 to the bytes buffer
	n := VarUInt(l + 1)
	if err := n.Marshal(b); err != nil {
		return fmt.Errorf("compact_array.length: %w", err)
	}

	if l == 0 {
		return nil
	}

	for i, v := range a.Values {
		// Get the underlying type of the value
		valType := reflect.TypeOf(v)
		if valType == nil {
			return fmt.Errorf("compact_array[%d]: nil value", i)
		}

		// Check for struct type
		// Marshal the struct using the datatypes.Marshal function
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			if err := Marshal(v, b); err != nil {
				return fmt.Errorf("compact_array[%d]: %w", i, err)
			}
			continue
		}

		// Must implement the MarshableDataType interface
		if !reflect.ValueOf(v).CanInterface() {
			return fmt.Errorf("field %s is not exportable", reflect.TypeOf(v).Name())
		}
		dt, ok := reflect.ValueOf(v).Interface().(MarshableDataType)
		if !ok {
			return fmt.Errorf("expected MarshableDataType, got %T", v)
		}
		if err := dt.Marshal(b); err != nil {
			return fmt.Errorf("compact_array[%d]: %w", i, err)
		}
	}

	return nil
}

func (a *CompactArray[T]) Unmarshal(p *Parser) error {
	// Unmarshal the length
	var length VarUInt
	if err := length.Unmarshal(p); err != nil {
		return fmt.Errorf("compact_array.length: %w", err)
	}

	// If length is 0, the array is nil
	if length == 0 {
		a.Values = nil
		return nil
	}

	// The encoded length is n+1, so we need to adjust
	length = length - 1
	a.Values = make([]T, length)

	// Create a sample value to get the type information
	var sampleVal T
	valType := reflect.TypeOf(sampleVal)
	valKind := valType.Kind()

	for i := range length {
		newValPtr := reflect.New(valType).Interface()

		if valKind == reflect.Struct {
			// For struct types, use the global Unmarshal function
			if err := Unmarshal(newValPtr, p); err != nil {
				return fmt.Errorf("compact_array[%d]: %w", i, err)
			}
			a.Values[i] = reflect.ValueOf(newValPtr).Elem().Interface().(T)
		} else {
			// For non-struct types, it must implement UnmarshableDataType
			dt, ok := newValPtr.(UnmarshableDataType)
			if !ok {
				return fmt.Errorf("compact_array[%d]: type %s does not implement UnmarshableDataType", i, valType.String())
			}

			if err := dt.Unmarshal(p); err != nil {
				return fmt.Errorf("compact_array[%d]: %w", i, err)
			}
			a.Values[i] = reflect.ValueOf(newValPtr).Elem().Interface().(T)
		}
	}

	return nil
}
