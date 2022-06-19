package service

import alfa "github.com/datafony/alfa"

func NewFooRouter(service FooService) *FooRouter {
	return &FooRouter{service: service}
}

type FooRouter struct {
	service FooService
}

func (f FooRouter) Routes() []alfa.Route {
	return []alfa.Route{
		{Method: alfa.Get, Handler: f.service.QueryPost},
		//{Method: alfa.Get, Handler: f.service.Ping, Description: "测试"},
		//{Method: alfa.Get, Handler: f.service.GetHtml},
		//{Method: alfa.Get, Handler: f.service.GetText},
		//{Method: alfa.Get, Handler: f.service.GetInt},
		//{Method: alfa.Get, Handler: f.service.GetInt32},
		//{Method: alfa.Get, Handler: f.service.GetDecimal},
		//{Method: alfa.Get, Handler: f.service.GetBool},
		//{Handler: f.service.AddPost},
		//{Handler: f.service.TestGetArray},
		//{Handler: f.service.TestPostArray},
		//{Handler: f.service.Ping2},
		//{Handler: f.service.PostShop},
	}
}
