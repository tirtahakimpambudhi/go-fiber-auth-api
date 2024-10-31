package reflecthelper

import (
	"fmt"
	"reflect"
	"strings"
)

// KeyValueToString converts the key-value pairs of a struct or map to a string.
func KeyValueToString(input interface{}) string {
	v := reflect.ValueOf(input)
	t := v.Type()

	var result strings.Builder

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldName := t.Field(i).Name
			fieldValue := v.Field(i).Interface()
			result.WriteString(fmt.Sprintf("%s: %v\n", fieldName, fieldValue))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			fieldName := key.Interface()
			fieldValue := v.MapIndex(key).Interface()
			result.WriteString(fmt.Sprintf("%v: %v\n", fieldName, fieldValue))
		}
	default:
		return "Invalid input: expected a struct or map"
	}

	return result.String()
}
