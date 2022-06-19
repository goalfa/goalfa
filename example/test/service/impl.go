package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type UserImplService struct {
}

func (u UserImplService) TestNestedInput(ctx context.Context, in A) (out A, err error) {
	out = in
	return
}

func (u UserImplService) TestDecimal(ctx context.Context) (out decimal.Decimal, err error) {
	out = decimal.NewFromFloat(3.14156)
	return
}

func (u UserImplService) Ping(ctx context.Context) (r string, err error) {
	r = "ok"
	return
}

func (u UserImplService) InnerError(ctx context.Context) (err error) {
	return errors.New("inner error")
}

func (u UserImplService) ValidateError(ctx context.Context) (err error) {
	return
}

func (u UserImplService) ForbiddenError(ctx context.Context) (err error) {
	return
}

func (u UserImplService) Wait(ctx context.Context) (err error) {
	time.Sleep(15 * time.Second)
	return
}

func (u UserImplService) TestGet(ctx context.Context, in TestStruct) (out TestStruct, err error) {
	out = in
	return
}

func (u UserImplService) TestGetArray(ctx context.Context, in TestStructs) (out TestStructs, err error) {
	out = in
	return
}

func (u UserImplService) TestPost(ctx context.Context, in TestStruct) (out TestStruct, err error) {
	out = in
	return
}

func (u UserImplService) TestPostArray(ctx context.Context, in TestStructs) (out TestStructs, err error) {
	out = in
	return
}

func (u UserImplService) TestDelete(ctx context.Context, in TestStruct) (out TestStruct, err error) {
	out = in
	return
}

func (u UserImplService) TestDeleteArray(ctx context.Context, in TestStructs) (out TestStructs, err error) {
	out = in
	return
}

func (u UserImplService) TestPut(ctx context.Context, in TestStruct) (out TestStruct, err error) {
	out = in
	return
}

func (u UserImplService) TestPutArray(ctx context.Context, in TestStructs) (out TestStructs, err error) {
	out = in
	return
}
