package foo_service

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_reflectService(t *testing.T) {
	//var s FooService
	//r1 := reflect.TypeOf(s)
	//t.Log("r1", r1.String())
	//s = new(FooMock)
	//r := reflect.ValueOf(s)
	//t.Log("s type:", r.Type())
	//for i := 0; i < r.NumMethod(); i++ {
	//	t.Log(r.Type().Method(i).Name)
	//	t.Log(r.Type().Method(i))
	//}
	
	sr := reflect.TypeOf((*FooService)(nil))
	fmt.Println(sr.Elem().NumMethod())
	//fmt.Println(sr.NumField())
}
