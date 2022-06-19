package service

import (
	"context"
	"github.com/shopspring/decimal"
)

type UserService interface {
	Ping(ctx context.Context) (r string, err error)
	Wait(ctx context.Context) (err error)
	InnerError(ctx context.Context) (err error)
	ValidateError(ctx context.Context) (err error)
	ForbiddenError(ctx context.Context) (err error)
	TestGet(ctx context.Context, in TestStruct) (out TestStruct, err error)
	TestGetArray(ctx context.Context, in TestStructs) (out TestStructs, err error)
	TestPost(ctx context.Context, in TestStruct) (out TestStruct, err error)
	TestPostArray(ctx context.Context, in TestStructs) (out TestStructs, err error)
	TestDelete(ctx context.Context, in TestStruct) (out TestStruct, err error)
	TestDeleteArray(ctx context.Context, in TestStructs) (out TestStructs, err error)
	TestPut(ctx context.Context, in TestStruct) (out TestStruct, err error)
	TestPutArray(ctx context.Context, in TestStructs) (out TestStructs, err error)
	TestDecimal(ctx context.Context) (out decimal.Decimal, err error)
	TestNestedInput(ctx context.Context, in A) (out A, err error)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TestStruct struct {
	String  string          `json:"string" validator:"required" label:"字符串"`
	Bool    bool            `json:"bool"`
	Int     int             `json:"int"`
	Int8    int8            `json:"int8"`
	Int16   int16           `json:"int16"`
	Int32   int32           `json:"int32"`
	Int64   int64           `json:"int64"`
	Uint    uint            `json:"uint"`
	Uint8   uint8           `json:"uint8"`
	Uint16  uint16          `json:"uint16"`
	Uint32  uint32          `json:"uint32"`
	Uint64  uint64          `json:"uint64"`
	Float32 float32         `json:"float32"`
	Float64 float32         `json:"float64"`
	Decimal decimal.Decimal `json:"decimal"`
	User    User            `json:"user"`
}

type TestStructs struct {
	Strings  []string          `json:"strings"`
	Bools    []bool            `json:"bools"`
	Ints     []int             `json:"ints"`
	Int8s    []int8            `json:"int8s"`
	Int16s   []int16           `json:"int16s"`
	Int32s   []int32           `json:"int32s"`
	Int64s   []int64           `json:"int64s"`
	Uints    []uint            `json:"uint"`
	Uint8s   []uint8           `json:"uint8s"`
	Uint16s  []uint16          `json:"uint16s"`
	Uint32s  []uint32          `json:"uint32s"`
	Uint64s  []uint64          `json:"uint64s"`
	Float32s []float32         `json:"float32s"`
	Float64s []float32         `json:"float64s"`
	Decimals []decimal.Decimal `json:"decimals"`
	Users    []User            `json:"users" label:"用户"`
}

type A struct {
	//A1 string `json:"a1"`
	//B1 B      `json:"b1"`
	BB []B    `json:"bb"`
	//C  C      `json:"c"`
}

type B struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
	C  C      `json:"c"`
}

type C struct {
	C1 string `json:"c1"`
}
