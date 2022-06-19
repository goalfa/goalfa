package service

import (
	"context"
	alfa "github.com/datafony/alfa"
	"github.com/shopspring/decimal"
)

type FooMock struct {
}

func (f FooMock) GetHtml(ctx context.Context) (out alfa.Html, err error) {
	out = "<h1>Hello world</h1>"
	return
}

func (f FooMock) GetText(ctx context.Context) (out alfa.Text, err error) {
	out = "plain text"
	return
}

func (f FooMock) GetDecimal(ctx context.Context) (out decimal.Decimal, err error) {
	out = decimal.NewFromFloat(3.14156)
	return
}

func (f FooMock) GetBool(ctx context.Context) (out bool, err error) {
	out = true
	return
}

func (f FooMock) GetInt(ctx context.Context) (out int, err error) {
	out = 1438
	return
}

func (f FooMock) GetInt32(ctx context.Context) (out *int, err error) {
	a := 123
	out = &a
	return
}

func (f FooMock) PostShop(ctx context.Context, in *Shop) (out *Shop, err error) {
	out = in
	return
}

func (f FooMock) Ping2(ctx context.Context) (out decimal.Decimal, err error) {
	out = decimal.NewFromFloat(3.14)
	return
}

func (f FooMock) TestGetArray(ctx context.Context) (out [][]string, err error) {
	out = [][]string{
		{"a", "n"},
	}
	return
}

func (f FooMock) TestPostArray(ctx context.Context) (out [][]Post, err error) {
	return
}

func (f FooMock) Ping(ctx context.Context) (out string, err error) {
	out = "ok"
	return
}

func (f FooMock) AddPost(ctx context.Context, in Post) (out Post, err error) {
	out = in
	return
}

func (f FooMock) QueryPost(ctx context.Context, in QueryPostIn) (out []Post, err error) {
	out = append(out, Post{
		Title:   "一篇文章",
		Content: "文章内容",
		Tags:    []string{"A", "B", "C"},
	})
	return
}
