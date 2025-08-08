package utils

import (
	"fmt"
	"reflect"
)

func UpdateField[T any](target *T, source *T) {
	if source != nil {
		*target = *source
	}
}

func UpdateFields(dest interface{}, src interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	if srcValue.Kind() != reflect.Ptr || destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both src and dest must be pointers")
	}

	srcElem := srcValue.Elem()
	destElem := destValue.Elem()

	if srcElem.Kind() != reflect.Struct || destElem.Kind() != reflect.Struct {
		return fmt.Errorf("both src and dest must be pointers to structs")
	}

	srcType := srcElem.Type()
	for i := 0; i < srcElem.NumField(); i++ {
		srcField := srcElem.Field(i)
		srcFieldType := srcType.Field(i)

		if !srcField.CanInterface() {
			continue
		}

		if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
			continue
		}

		destField := destElem.FieldByName(srcFieldType.Name)
		if !destField.IsValid() || !destField.CanSet() {
			continue
		}

		if srcField.Kind() == reflect.Ptr {
			destField.Set(srcField.Elem())
		} else {
			destField.Set(srcField)
		}
	}

	return nil
}

func ApplyPartialUpdate[T any, P any](existing *T, payload *P) error {
	return UpdateFields(existing, payload)
}
