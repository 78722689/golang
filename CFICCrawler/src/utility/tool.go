package utility

import (
	"reflect"
	"fmt"
	"os"
)


// To check if the list contains the elem
func Contains(list interface{}, elem interface{}) bool {
	value := reflect.ValueOf(list)
	if value.Kind() != reflect.Slice {
		fmt.Fprintf(os.Stderr, "Input type is not an array or slice type: %v, kind:%s", value, value.Kind())
		return false
	}

	for i:=0; i<value.Len();i++ {
		if value.Index(i).Interface() == elem.(string) {
			return true
		}
	}

	return false
}

// To get map keys
func Keys(i interface{}) (keys []interface{}) {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Map {
		fmt.Fprintf(os.Stderr, "Input type is not a map type: %v", v)
		return nil
	}

	for _,key := range v.MapKeys() {
		keys = append(keys, key.Interface())
	}

	return keys
}