package utils

import "reflect"

func ValueElem(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Ptr {
			return v
		}
		v = v.Elem()
	}
}

func TypeElem(v reflect.Type) reflect.Type {
	for {
		if v.Kind() != reflect.Ptr {
			return v
		}
		v = v.Elem()
	}
}
