package request

import (
	"fmt"
	"reflect"
)

// getTagPart splits parse the first subtag from the tag.
// The tags are expected to be in the format "kafka:<tag1>:<tag2>:<tag3>..".
// It returns the first subtag and the rest of the tag.
func getTagPart(tag string) (string, string) {
	for i, c := range tag {
		if c == ':' {
			return tag[:i], tag[i+1:]
		}
	}

	return tag, ""
}

// Decodes a value based on the kafka tags in it's fields.
func decode(r *RequestHeader, tag string, v reflect.Value, path string) error {
	tag, nextTag := getTagPart(tag)

	switch tag {
	case "compact_array":
		if nextTag == "" {
			return fmt.Errorf("%s: 'compact_array' tag must be accompanied by a element type", path)
		}
		return decodeCompactArray(r, nextTag, v, path)

	case "struct":
		return decodeStruct(r, v, path, false)

	case "inline_struct":
		return decodeStruct(r, v, path, true)

	default:
		return decodePrimitive(r, tag, v, path)
	}
}

// Decodes a compact array based on the kafka tags in it's fields.
// The tag is expected to be the tag of the element type.
func decodeCompactArray(r *RequestHeader, tag string, v reflect.Value, path string) error {
	length, err := r.ReadVarUInt()
	if err != nil {
		return fmt.Errorf("%s.length: %w", path, err)
	}

	// Handle the case of zero length: null array
	if length == 0 {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		return nil
	}
	length--

	slice := reflect.MakeSlice(v.Type(), int(length), int(length))
	for i := range length {
		elemPath := fmt.Sprintf("%s[%d]", path, i)
		elemPtr := reflect.New(v.Type().Elem())
		if err := decode(r, tag, elemPtr, elemPath); err != nil {
			return err
		}
		slice.Index(int(i)).Set(elemPtr.Elem())
	}
	v.Set(slice)

	return nil
}

// Decodes a struct based on the kafka tags in it's fields.
// A regular struct is expected to have empty tagged fields at it's end, whereas
// an inline struct has no such requirement.
func decodeStruct(r *RequestHeader, v reflect.Value, path string, isInlineStruct bool) error {
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("%s: must be a non-nil pointer to struct", path)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%s: must point to a struct", path)
	}
	t := v.Type()

	for i := range t.NumField() {
		f := t.Field(i)
		subPath := fmt.Sprintf("%s.%s", path, f.Name)

		fieldVal := v.Field(i)
		tag := f.Tag.Get("kafka")

		if err := decode(r, tag, fieldVal, subPath); err != nil {
			return err
		}
	}

	if !isInlineStruct {
		if err := r.ReadZeroTaggedFieldArray(); err != nil {
			return fmt.Errorf("%s.taggedFields: %w", path, err)
		}
	}

	return nil
}

// Decodes a primitive type based on the tag.
func decodePrimitive(r *RequestHeader, tag string, v reflect.Value, path string) error {
	switch tag {

	case "int8":
		val, err := r.ReadInt8()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetInt(int64(val))

	case "int16":
		val, err := r.ReadInt16()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetInt(int64(val))

	case "int32":
		val, err := r.ReadInt32()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetInt(int64(val))

	case "int64":
		val, err := r.ReadInt64()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetInt(val)

	case "varint":
		val, err := r.ReadVarInt()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetInt(int64(val))

	case "varuint":
		val, err := r.ReadVarUInt()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetUint(uint64(val))

	// String types
	case "string":
		val, err := r.ReadString()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetString(val)

	case "compact_string":
		val, err := r.ReadCompactString()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetString(val)

	case "nullable_string":
		val, err := r.ReadNullableString()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if val == nil {
			v.Set(reflect.Zero(v.Type()))
		} else {
			v.Set(reflect.ValueOf(val))
		}

	case "uuid":
		val, err := r.ReadUUID()
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		v.SetBytes(val)

	default:
		return fmt.Errorf("%s: unknown tag %q", path, tag)
	}

	return nil
}
