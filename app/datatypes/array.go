package datatypes

import (
	"bytes"
	"fmt"
	"reflect"
)

// CompactArray stores the length of the array as a VarUInt, and the values as a slice of T.
// If the array is nil, the length is 0. Otherwise, the length is the number of elements in the array + 1.
type CompactArray[T any] struct {
	Values []T
}

// Array stores the length of the array as an Int32, and the values as a slice of T.
type Array[T any] struct {
	Values []T
}

// VarIntArray stores the length of the array as a VarInt, and the values as a slice of T.
type VarIntArray[T any] struct {
	Values []T
}

func (a *CompactArray[T]) Append(v *T) {
	a.Values = append(a.Values, *v)
}

// marshalArrayUtil is a utility function to marshal arrays of various types.
// It handles both nil and non-nil cases, and uses the provided length
// and nil length values to determine how to marshal the data.
func marshalArrayUtil[T any](vals []T, nilLen MarshableDataType, length MarshableDataType, b *bytes.Buffer, typeName string) error {
	if vals == nil {
		if err := nilLen.Marshal(b); err != nil {
			return fmt.Errorf("%s.nil_array: %w", typeName, err)
		}
		return nil
	}

	if err := length.Marshal(b); err != nil {
		return fmt.Errorf("%s.length: %w", typeName, err)
	}

	for i, v := range vals {
		// Get the underlying type of the value
		valType := reflect.TypeOf(v)
		if valType == nil {
			return fmt.Errorf("%s[%d]: nil value", typeName, i)
		}

		// Check for struct type
		if valType.Kind() == reflect.Struct {
			if err := Marshal(v, b); err != nil {
				return fmt.Errorf("%s[%d]: %w", typeName, i, err)
			}
			continue
		}

		dt, ok := reflect.ValueOf(v).Interface().(MarshableDataType)
		if !ok {
			return fmt.Errorf("%s[%d]: expected MarshableDataType, got %T", typeName, i, v)
		}
		if err := dt.Marshal(b); err != nil {
			return fmt.Errorf("%s[%d]: %w", typeName, i, err)
		}
	}
	return nil
}

func (a CompactArray[T]) Marshal(b *bytes.Buffer) error {
	nilLen := VarUInt(0)
	length := VarUInt(len(a.Values) + 1)
	if err := marshalArrayUtil(a.Values, nilLen, length, b, "compact_array"); err != nil {
		return err
	}
	return nil
}

func (a Array[T]) Marshal(b *bytes.Buffer) error {
	nilLen := Int32(0)
	length := Int32(len(a.Values))
	if err := marshalArrayUtil(a.Values, nilLen, length, b, "array"); err != nil {
		return err
	}
	return nil
}

func (a VarIntArray[T]) Marshal(b *bytes.Buffer) error {
	nilLen := VarInt(0)
	length := VarInt(len(a.Values))
	if err := marshalArrayUtil(a.Values, nilLen, length, b, "varint_array"); err != nil {
		return err
	}
	return nil
}

func unmarshalArrayUtil[T any](vals []T, p *Parser, typeName string) error {
	// Create a sample value to get the type information
	var sampleVal T
	valType := reflect.TypeOf(sampleVal)
	valKind := valType.Kind()

	for i := range vals {
		newValPtr := reflect.New(valType).Interface()

		if valKind == reflect.Struct {
			// For struct types, use the global Unmarshal function
			if err := Unmarshal(newValPtr, p); err != nil {
				return fmt.Errorf("%s[%d]: %w", typeName, i, err)
			}
			vals[i] = reflect.ValueOf(newValPtr).Elem().Interface().(T)
		} else {
			// For non-struct types, it must implement UnmarshableDataType
			dt, ok := newValPtr.(UnmarshableDataType)
			if !ok {
				return fmt.Errorf("%s[%d]: type %s does not implement UnmarshableDataType", typeName, i, valKind.String())
			}

			if err := dt.Unmarshal(p); err != nil {
				return fmt.Errorf("%s[%d]: %w", typeName, i, err)
			}
			vals[i] = reflect.ValueOf(newValPtr).Elem().Interface().(T)
		}
	}

	return nil
}

func (a *CompactArray[T]) Unmarshal(p *Parser) error {
	var length VarUInt
	if err := length.Unmarshal(p); err != nil {
		return fmt.Errorf("compact_array.length: %w", err)
	}

	if length == 0 {
		a.Values = nil
		return nil
	}

	length = length - 1
	a.Values = make([]T, length)

	if err := unmarshalArrayUtil(a.Values, p, "compact_array"); err != nil {
		return fmt.Errorf("compact_array: %w", err)
	}

	return nil
}

func (a *Array[T]) Unmarshal(p *Parser) error {
	var length Int32
	if err := length.Unmarshal(p); err != nil {
		return fmt.Errorf("array.length: %w", err)
	}

	if length == -1 {
		a.Values = nil
		return nil
	}

	a.Values = make([]T, length)
	if err := unmarshalArrayUtil(a.Values, p, "array"); err != nil {
		return fmt.Errorf("array: %w", err)
	}

	return nil
}

func (a *VarIntArray[T]) Unmarshal(p *Parser) error {
	var length VarInt
	if err := length.Unmarshal(p); err != nil {
		return fmt.Errorf("varint_array.length: %w", err)
	}

	a.Values = make([]T, length)
	if err := unmarshalArrayUtil(a.Values, p, "varint_array"); err != nil {
		return fmt.Errorf("varint_array: %w", err)
	}

	return nil
}
