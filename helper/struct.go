package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Get a struct field by its tag value
func getFieldName[T any](tag string, key string) (string, error) {
	var t T

	fieldType := reflect.TypeOf(t)

	// Prevention check for non struct types
	if fieldType.Kind() != reflect.Struct {
		return "", fmt.Errorf("the type must be struct")
	}

	// Iterating over the field by index to find the tag value
	for i := 0; i < fieldType.NumField(); i++ {
		field := fieldType.Field(i)

		// Get the tag value
		valueTag := strings.Split(
			field.Tag.Get(key),
			",",
		)[0]

		if valueTag == tag {
			return field.Name, nil
		}
	}

	return "", fmt.Errorf("the field doesn't exist")
}

// Get a struct field from its JSON tag value
func GetFieldNameFromJSON[T any](tag string) (string, error) {
	return getFieldName[T](tag, "json")
}

// Get ptr value from interface{}
func GetValuePtr(valueRaw any) any {
	p := reflect.New(reflect.TypeOf(valueRaw))

	p.Elem().Set(reflect.ValueOf(valueRaw))

	return p.Interface()
}

// Return the field as a `reflect.Value` if it is mutable and valid
func getFieldValue[T any](
	elem *reflect.Value,
	nameJSON string,
) *reflect.Value {
	fieldName, err := GetFieldNameFromJSON[T](nameJSON)

	if err != nil {
		return nil
	}

	field := elem.FieldByName(fieldName)

	if !field.IsValid() || !field.CanSet() {
		return nil
	}

	return &field
}

// Assign a value to `field`
func setFieldValue(field *reflect.Value, valueRaw any) bool {
	var value any

	switch field.Kind() {

	case reflect.Array | reflect.Slice:
		return false

	case reflect.Ptr:
		valuePtr := GetValuePtr(valueRaw)

		if valuePtr == nil {
			return false
		}

		value = valuePtr

	default:
		value = valueRaw
	}

	field.Set(reflect.ValueOf(value))

	return true
}

// Returns a struct `T` from a Terraform schema keys `schema` and `values`
func FromSchemas[T any](
	schemas map[string]*schema.Schema,
	d *schema.ResourceData,
) T {
	var t T

	elem := reflect.ValueOf(&t).Elem()

	for key := range schemas {
		valueRaw, ok := d.GetOk(key)

		if valueRaw == nil || !ok {
			continue
		}

		field := getFieldValue[T](&elem, key)

		if field == nil {
			continue
		}

		setFieldValue(field, valueRaw)
	}

	return elem.Interface().(T)
}

// Returns a struct from a specific map
func FromMap[T any](m map[string]any) T {
	var t T

	elem := reflect.ValueOf(&t).Elem()

	for key, value := range m {
		field := getFieldValue[T](&elem, key)

		if field == nil {
			continue
		}

		setFieldValue(field, value)
	}

	return elem.Interface().(T)
}

// Only for `schema.TypeSet` with `schema.Resource`
func FromTypeSetCallback[T any](
	set []any,
	f func(map[string]any) (*T, error),
) []T {
	var t []T

	for _, mapRaw := range set {
		_map := mapRaw.(map[string]any)
		value, err := f(_map)

		if err != nil {
			continue
		}

		t = append(t, *value)
	}

	return t
}

// Only for `schema.TypeSet` with `schema.Resource`
func FromTypeSetDefault[T any](set []any) []T {
	return FromTypeSetCallback(
		set, func(m map[string]any) (*T, error) {
			ret := FromMap[T](m)

			return &ret, nil
		},
	)
}

// Return a slice from a `schema.TypeList` or `schema.TypeSet`
func FromSlice[T any](slice []any) ([]T, error) {
	var t []T

	for _, value := range slice {
		valueAssertion, ok := value.(T)

		if !ok {
			return []T{}, fmt.Errorf("unable to assert")
		}

		t = append(t, valueAssertion)
	}

	return t, nil
}
