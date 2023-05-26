package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Get a struct field as reflect.Value
func getAttr(obj any, fieldName string) (reflect.Value, error) {
	s := reflect.ValueOf(obj)
	field := s.FieldByName(fieldName)

	if !field.IsValid() {
		return s, fmt.Errorf("invalid field")
	}

	return field, nil
}

// Returns if a value is a struct
func isStruct(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Struct
}

// Update a value in a Terraform ResourceData
func setResourceField(
	k string,
	v any,
	d *schema.ResourceData,
) bool {
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		return false
	}

	dValue := d.Get(k)

	if dValue == nil {
		return false
	}

	d.Set(k, v)

	return true
}

// Convert a struct to map
//
// Faster way to use json routines
func StructToMap(s any) (map[string]any, error) {
	m := make(map[string]any)

	if !isStruct(s) {
		return m, fmt.Errorf("must be a struct")
	}

	fields := reflect.VisibleFields(reflect.TypeOf(s))

	for _, field := range fields {
		fieldValue, err := getAttr(s, field.Name)

		if err != nil {
			continue
		}

		tag := strings.Split(
			field.Tag.Get("json"),
			",",
		)

		if len(tag) < 1 {
			return m, fmt.Errorf("must have a JSON tag")
		}

		m[tag[0]] = fieldValue.Interface()
	}

	return m, nil
}

// Update values in a Terraform ResourceData from `s`
func ResourceFromStruct(s any, d *schema.ResourceData) error {
	m, err := StructToMap(s)

	if err != nil {
		return err
	}

	for k, v := range m {
		setResourceField(k, v, d)
	}

	return nil
}

// Get a new composed set from a struct
func NewSetCallback[T any](
	structs *[]T,
	f func(s any) (map[string]any, error),
) ([]map[string]any, error) {
	items := []map[string]any{}

	for _, s := range *structs {
		_map, err := f(s)

		if err != nil {
			return nil, err
		}

		items = append(items, _map)
	}

	return items, nil
}

// Returns a new composed set with a default function
func NewSetDefault[T any](structs *[]T) ([]map[string]any, error) {
	return NewSetCallback(structs, StructToMap)
}
