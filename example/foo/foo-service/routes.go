package foo_service

import "github.com/koyeo/goalfa"

func NewFooRouter(service FooService) *FooRouter {
	return &FooRouter{service: service}
}

type FooRouter struct {
	service FooService
}

func (f FooRouter) Routes() []goalfa.Route {
	return []goalfa.Route{
		{
			Service: goalfa.Service{
				Interface: (*FooService)(nil),
				Implement: f.service,
			},
			Routes: []goalfa.Route{
				{Method: goalfa.Get, Path: "/download", Handler: f.service.Ping},
				//{Method: goalfa.Get, Handler: f.service.GetUserByIdStruct},
				{Method: goalfa.Get, Params: []goalfa.Param{{Field: "user_id"}}, Handler: f.service.GetUserById},
				{Method: goalfa.Get, Params: []goalfa.Param{{Field: "group_id"}}, Handler: f.service.GetUserById},
				{Method: goalfa.Get, Params: []goalfa.Param{{Field: "id"}, {Field: "filter"}}, Handler: f.service.GetUserById},
			},
		},
	}
}
