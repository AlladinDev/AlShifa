package utils

import (
	"net/url"
	"reflect"
	"strconv"
)

// UnmarshalFormValues maps url.Values into a struct pointer using `form` tags
// or the struct field name if no tag is present.
func UnmarshalFormValues(form url.Values, target any) {
	if target == nil {
		return
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		key := structField.Tag.Get("Form")
		if key == "-" {
			continue
		}
		if key == "" {
			key = structField.Name
		}

		values, ok := form[key]
		if !ok || len(values) == 0 {
			continue
		}

		// Handle pointer fields
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		switch field.Kind() {

		case reflect.String:
			field.SetString(values[0])

		case reflect.Bool:
			if v, err := strconv.ParseBool(values[0]); err == nil {
				field.SetBool(v)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(values[0], 10, field.Type().Bits()); err == nil {
				field.SetInt(v)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(values[0], 10, field.Type().Bits()); err == nil {
				field.SetUint(v)
			}

		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(values[0], field.Type().Bits()); err == nil {
				field.SetFloat(v)
			}

		case reflect.Slice:
			elemKind := field.Type().Elem().Kind()

			switch elemKind {

			case reflect.String:
				field.Set(reflect.ValueOf(values))

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				slice := reflect.MakeSlice(field.Type(), 0, len(values))
				for _, s := range values {
					if v, err := strconv.ParseInt(s, 10, field.Type().Elem().Bits()); err == nil {
						slice = reflect.Append(slice, reflect.ValueOf(v).Convert(field.Type().Elem()))
					}
				}
				field.Set(slice)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				slice := reflect.MakeSlice(field.Type(), 0, len(values))
				for _, s := range values {
					if v, err := strconv.ParseUint(s, 10, field.Type().Elem().Bits()); err == nil {
						slice = reflect.Append(slice, reflect.ValueOf(v).Convert(field.Type().Elem()))
					}
				}
				field.Set(slice)

			case reflect.Float32, reflect.Float64:
				slice := reflect.MakeSlice(field.Type(), 0, len(values))
				for _, s := range values {
					if v, err := strconv.ParseFloat(s, field.Type().Elem().Bits()); err == nil {
						slice = reflect.Append(slice, reflect.ValueOf(v).Convert(field.Type().Elem()))
					}
				}
				field.Set(slice)

			case reflect.Bool:
				slice := reflect.MakeSlice(field.Type(), 0, len(values))
				for _, s := range values {
					if v, err := strconv.ParseBool(s); err == nil {
						slice = reflect.Append(slice, reflect.ValueOf(v))
					}
				}
				field.Set(slice)
			}
		}
	}
}
