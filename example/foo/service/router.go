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
		{Handler: f.service.Ping},
		{Handler: f.service.SaveUser},
		//{Method: alfa.Get, Path: "/GetUserById/:id", Handler: f.service.GetUserById},
		{Method: alfa.Get, Handler: f.service.GetUserByIdStruct},
		{Handler: f.service.SaveArticle},
	}
}
