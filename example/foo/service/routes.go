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
		{
			Prefix:  "/api",
			Service: f.service,
			Routes: []alfa.Route{
				{Method: alfa.Get, Handler: f.service.Ping},
				{Method: alfa.Get, Handler: f.service.GetUserByIdStruct},
				{Method: alfa.Get, Params: []alfa.Param{{Field: "id", Description: "客户ID"}}, Handler: f.service.GetUserById},
			},
		},
	}
}
