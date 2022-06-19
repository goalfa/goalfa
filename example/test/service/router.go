package service

func NewUserRouter(service UserService) *UserRouter {
	return &UserRouter{service: service}
}

type UserRouter struct {
	service UserService
}

func (p UserRouter) Routes() []buck.Route {
	return []buck.Route{
		{
			Prefix: "/api",
			Children: []buck.Route{
				{
					Prefix: "",
					Children: []buck.Route{
						{Method: buck.Get, Path: "/ping", Handler: p.service.Ping, Description: "测试项目是否连通"},
						{Method: buck.Get, Handler: p.service.Wait},
						{Method: buck.Get, Path: "/inner-error", Handler: p.service.InnerError},
						{Method: buck.Get, Handler: p.service.ValidateError},
						{Method: buck.Get, Handler: p.service.ForbiddenError},
					},
				},
				{
					Prefix: "/guard",
					Children: []buck.Route{
						{Handler: p.service.TestPost},
						{Handler: p.service.TestPostArray},
						{Method: buck.Get, Handler: p.service.TestGet},
						{Method: buck.Get, Handler: p.service.TestGetArray},
						{Method: buck.Put, Handler: p.service.TestPut},
						{Method: buck.Put, Handler: p.service.TestPutArray},
						{Method: buck.Delete, Handler: p.service.TestDelete},
						{Method: buck.Delete, Handler: p.service.TestDeleteArray},
						{Method: buck.Get, Handler: p.service.TestDecimal},
						{Method: buck.Post, Handler: p.service.TestNestedInput},
					},
				},
			},
		},
	}
}
