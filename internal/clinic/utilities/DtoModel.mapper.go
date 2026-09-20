package utilities

import (
	"reflect"
)

func DtoToModelMapper(structOne interface{}, structTwo interface{}) {
	src := reflect.ValueOf(structOne)
	dst := reflect.ValueOf(structTwo)

	// Dereference pointers
	if src.Kind() == reflect.Ptr {
		if src.IsNil() {
			return
		}
		src = src.Elem()
	}

	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		// Destination must be a non-nil pointer
		return
	}

	dst = dst.Elem()

	// Both must be structs
	if src.Kind() != reflect.Struct || dst.Kind() != reflect.Struct {
		return
	}

	srcType := src.Type()

	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		srcFieldType := srcType.Field(i)

		// Ignore unexported fields
		if !srcField.CanInterface() {
			continue
		}

		// Find matching field in destination by name
		dstField := dst.FieldByName(srcFieldType.Name)

		// Field must exist and be settable
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		// Types must be assignable
		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
		}
	}
}
