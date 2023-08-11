package autowire

import (
	"fmt"
	"reflect"
)

type Type = reflect.Type

// TypeOf 获取反射类型
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func getTypeName[T any]() string {
	return getTypeNameT(TypeOf[T]())
}

// 如果是指针类型，持续解引用，直到得到底层值类型
func getTypeNameT(typ reflect.Type) string {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.PkgPath() == "" {
		return typ.Name()
	}

	return fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
}

func SetValue[T any](receiver *T, v any) {
	reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(v))
}
