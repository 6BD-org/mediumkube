package utils

import (
	"reflect"
)

// Contains whether a slice contains an element
func Contains(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// ContainsT If slice contains elem
func ContainsT(slice interface{}, elem interface{}) bool {
	elemVal := reflect.ValueOf(elem)
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		sliceVal := reflect.ValueOf(slice)
		for i := 0; i < sliceVal.Len(); i++ {
			if reflect.DeepEqual(
				sliceVal.Index(i),
				elemVal,
			) {
				return true
			}
		}
	}

	return false
}
