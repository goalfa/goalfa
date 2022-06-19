package service

import (
	"context"
	alfa "github.com/datafony/alfa"
	"github.com/shopspring/decimal"
)

type FooService interface {
	QueryPost(ctx context.Context, in QueryPostIn) (out []Post, err error)
	Ping(ctx context.Context) (out string, err error)
	GetHtml(ctx context.Context) (out alfa.Html, err error)
	GetText(ctx context.Context) (out alfa.Text, err error)
	GetInt(ctx context.Context) (out int, err error)
	GetInt32(ctx context.Context) (out *int, err error)
	GetDecimal(ctx context.Context) (out decimal.Decimal, err error)
	GetBool(ctx context.Context) (out bool, err error)
	Ping2(ctx context.Context) (out decimal.Decimal, err error)
	AddPost(ctx context.Context, in Post) (out Post, err error)
	TestGetArray(ctx context.Context) (out [][]string, err error)
	TestPostArray(ctx context.Context) (out [][]Post, err error)
	PostShop(ctx context.Context, shop *Shop) (out *Shop, err error)
}

type CID struct {
	value string
}

func (p CID) String() string {
	return p.value
}

type Post struct {
	Title   string   `json:"title,omitempty" label:"标题" validator:"required" description:""`
	Content string   `json:"content,omitempty" label:"内容"`
	Tags    []string `json:"tags,omitempty" label:"标签"`
}

type QueryPostIn struct {
	Page  int    `json:"page"`
	Limit string `json:"limit"`
	Order *Order `json:"order"`
}

type Order struct {
	Field     string
	Direction string
}

type Shop struct {
	Name      string     `json:"name,omitempty" label:"店铺名称"`
	Manager   *Member    `json:"manager,omitempty" label:"管理员"`
	Employees []*Member  `json:"employees,omitempty" label:"其他雇员"`
	Products  []*Product `json:"products,omitempty" label:"商品列表"`
}

type Member struct {
	Name   string `json:"name,omitempty" label:"雇员姓名"`
	Mobile string `json:"mobile,omitempty" label:"雇员电话"`
}

type Product struct {
	Name string  `json:"name,omitempty" label:"商品名称"`
	Spec []*Spec `json:"spec,omitempty" label:"商品规格"`
}

type Spec struct {
	Title     string          `json:"title,omitempty" label:"规格标题" json:"title,omitempty"`
	Inventory int             `json:"inventory,omitempty" label:"规格库存" json:"inventory,omitempty"`
	Price     decimal.Decimal `label:"规格价格" json:"price"`
}
