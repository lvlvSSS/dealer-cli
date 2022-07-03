package converter

import (
	"reflect"
	"unsafe"
)

func BytesToString(b []byte) string {
	sliceheader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceheader.Data,
		Len:  sliceheader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

/*
	StringToBytes is to convert string to slice with zero-copy.
*/
func StringToBytes(s string) []byte {
	stringheader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh := reflect.SliceHeader{
		Data: stringheader.Data,
		Len:  stringheader.Len,
		Cap:  stringheader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&sh))
}
