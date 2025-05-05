package response

import (
	"fmt"
	"reflect"
)

// getTagPart splits the kafka tag into first tag and remaining tags
func getTagPart(tag string) (string, string) {
	for i, c := range tag {
		if c == ':' {
			return tag[:i], tag[i+1:]
		}
	}
	return tag, ""
}

func encode(r *Response, tag string, v reflect.Value, path string) error {
	tag, nextTag := getTagPart(tag)

	switch tag {
	case "compact_array":
		if nextTag == "" {
			return fmt.Errorf("%s: 'compact_array' tag requires element type", path)
		}
		return encodeCompactArray(r, nextTag, v, path)

	case "struct":
		return encodeStruct(r, v, path, false)

	case "inline_struct":
		return encodeStruct(r, v, path, true)

	default:
		return encodePrimitive(r, tag, v, path)
	}
}

func encodeStruct(r *Response, v reflect.Value, path string, isInline bool) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%s: expected struct", path)
	}
	t := v.Type()

	for i := range v.NumField() {
		fieldVal := v.Field(i)
		field := t.Field(i)

		tag := field.Tag.Get("kafka")
		subPath := fmt.Sprintf("%s.%s", path, field.Name)

		if err := encode(r, tag, fieldVal, subPath); err != nil {
			return err
		}
	}

	if !isInline {
		r.WriteEmptyTaggedFields()
	}

	return nil
}

func encodeCompactArray(r *Response, tag string, v reflect.Value, path string) error {
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("%s: expected slice", path)
	}

	r.WriteCompactArrayLength(v.Len())
	for i := range v.Len() {
		elem := v.Index(i)
		elemPath := fmt.Sprintf("%s[%d]", path, i)
		if err := encode(r, tag, elem, elemPath); err != nil {
			return err
		}
	}
	return nil
}

func encodePrimitive(r *Response, tag string, v reflect.Value, path string) error {
	switch tag {
	case "int8":
		r.WriteInt8(int8(v.Int()))

	case "int16":
		r.WriteInt16(int16(v.Int()))

	case "int32":
		r.WriteInt32(int32(v.Int()))

	case "int64":
		r.WriteInt64(v.Int())

	case "varint":
		r.WriteVarInt(int(v.Int()))

	case "varuint":
		r.WriteVarUInt(uint(v.Uint()))

	case "string":
		r.WriteString(v.String())

	case "compact_string":
		r.WriteCompactString(v.String())

	case "nullable_string":
		if v.IsNil() {
			r.WriteNullableString(nil)
		} else {
			str := v.Elem().String()
			r.WriteNullableString(&str)
		}

	case "uuid":
		r.WriteUUID(v.Bytes())

	default:
		return fmt.Errorf("%s: unknown tag %q", path, tag)
	}
	return nil
}
