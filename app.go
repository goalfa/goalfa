package goalfa

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_dump"
	"github.com/koyeo/goalfa/binding"
	"github.com/koyeo/goalfa/exporter"
	"github.com/koyeo/goalfa/logger"
	"github.com/shopspring/decimal"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func New() *App {
	return &App{
		routeTable: &RouteTable{},
		handlers:   map[string]int{},
		routes:     map[string]string{},
	}
}

type App struct {
	version    string
	routers    []Router
	engine     *gin.Engine
	routeTable *RouteTable
	exporter   *exporter.Exporter
	methods    []*exporter.Method
	basics     *exporter.BasicTypes
	models     *exporter.Structs
	mode       Mode
	registers  []Route           // 最终注册的路由列表
	handlers   map[string]int    // 用以处理路由方法重复，并以最新值替换旧值
	routes     map[string]string // 用以检查路由（Method@Path） 是否重复
}

func (p *App) SetVersion(version string) {
	p.version = version
}

func (p *App) AddRouter(router ...Router) {
	p.routers = append(p.routers, router...)
}

func (p *App) SetEngine(engine *gin.Engine) {
	p.engine = engine
}

func (p *App) SetExporter(addr string, options *exporter.Settings) {
	basicTypes := []exporter.BasicType{
		{
			Elem: decimal.Decimal{},
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
		{
			Elem: time.Time{},
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
		{
			Elem: time.Duration(0),
			Mapping: map[string]exporter.Library{
				"ts": {Type: "number"},
			},
		},
		{
			Elem: Html(""),
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
		{
			Elem: Text(""),
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
	}
	if options == nil {
		options = new(exporter.Settings)
	}
	options.BasicTypes = append(basicTypes, options.BasicTypes...)
	p.exporter = exporter.NewExporter(addr, options)
}

func (p *App) Run(addr string) {
	if p.engine == nil {
		p.engine = gin.Default()
	}
	
	var err error
	
	for _, router := range p.routers {
		var routes []Route
		routes, err = p.prepareRoutes(router.Routes())
		if err != nil {
			logger.Error(err)
			return
		}
		p.addRegisterRoute(routes...)
	}
	
	_dump.Json(p.registers)
	//err = p.registerRoutes(p.engine, "", p.registers)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//
	//if p.exporter != nil {
	//	p.exporter.Init(p.version, p.methods, p.models)
	//	p.exporter.Run()
	//}
	//
	//err = p.engine.Run(addr)
	//if err != nil {
	//	panic(err)
	//	return
	//}
}

// 添加待注册的路由
// ⚠️ 如果自定义方法路由有被覆盖的情况，则输出警告
func (p *App) addRegisterRoute(routes ...Route) {
	for _, v := range routes {
		//if v.Handler != nil {
		//	p.handlers[]
		//}
		p.registers = append(p.registers, v)
	}
	
}

// 检查参数是否为 error 类型
func (p *App) isError(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

// 检查参数是否为 context 类型
func (p *App) isContext(v reflect.Type) bool {
	if v.Name() == "Context" && v.PkgPath() == "context" {
		return true
	}
	return false
}

// 检查参数是否接受的路由 Handler 格式
func (p *App) isHandler(t reflect.Type) error {
	if t.Kind() != reflect.Func {
		return fmt.Errorf("handler expect type func")
	}
	//if t.NumIn() != 1 && t.NumIn() != 2 {
	//	return fmt.Errorf("max 2 input parameters expected")
	//}
	// TODO Check Params Bind 参数
	//if t.NumIn() == 2 {
	//	in := t.In(1)
	//	for {
	//		if in.Kind() != reflect.Ptr {
	//			break
	//		}
	//		in = in.Elem()
	//	}
	//	if in.Kind() != reflect.Struct {
	//		return fmt.Errorf("input parameter only acept struct")
	//	}
	//}
	if !p.isContext(t.In(0)) {
		return fmt.Errorf("第一个入参期望 context.Context 类型，当前类型：%s", t.In(0).String())
	}
	if t.NumOut() != 1 && t.NumOut() != 2 {
		return fmt.Errorf("max 2 output parameters expected")
	}
	if !p.isError(t.Out(t.NumOut() - 1)) {
		return fmt.Errorf("last output parameter expect type error")
	}
	return nil
}

// 反射路由 Handler, 并检查是否为可接受的格式
func (p *App) parseHandler(handler interface{}) (v reflect.Value, err error) {
	v = reflect.ValueOf(handler)
	if err = p.isHandler(v.Type()); err != nil {
		err = fmt.Errorf("路由方法格式错误: %s: %s", v.Type(), err)
		return
	}
	return
}

// 预处理路由，反射路由处理器，并检查类型
func (p *App) prepareRoutes(routes []Route) (out []Route, err error) {
	out = make([]Route, len(routes))
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		route.Prefix = strings.TrimSpace(route.Prefix)
		route.Method = strings.TrimSpace(route.Method)
		route.Path = strings.TrimSpace(route.Path)
		
		// parse service and register exported methods
		if route.Service.Instance != nil {
			var serviceRoutes []Route
			serviceRoutes, err = p.prepareService(route.Service)
			if err != nil {
				return
			}
			route.Routes = append(serviceRoutes, route.Routes...)
		}
		
		if route.Handler != nil {
			route.handler, err = p.parseHandler(route.Handler)
			if err != nil {
				err = wrapHandlerError(route.Handler, err)
				return
			}
			route.handlerInfo = p.parseHandlerInfoValue(route.handler)
			if route.Method == "" {
				route.Method = http.MethodPost
			}
			if route.Path == "" {
				route.Path = fmt.Sprintf("/%s", route.handlerInfo.Name)
			}
		}
		
		route.Routes, err = p.prepareRoutes(route.Routes)
		if err != nil {
			return
		}
		
		out[i] = route
	}
	return
}

func methodInfo(method reflect.Value) *runtime.Func {
	return runtime.FuncForPC(method.Pointer())
}

func wrapHandlerError(handler interface{}, err error) error {
	v := reflect.ValueOf(handler)
	f := runtime.FuncForPC(v.Pointer())
	return fmt.Errorf("%s 路由解析错误: %s", f.Name(), err)
}

// 注册服务方法
func (p *App) prepareService(service Service) (out []Route, err error) {
	
	if service.Instance == nil {
		err = fmt.Errorf("服务绑定实例 Instance 为 nil")
		return
	}
	
	sv := reflect.ValueOf(service.Instance)
	if sv.NumMethod() == 0 {
		logger.Warn(fmt.Sprintf("服务 %s 没有可用的导出方法", sv.String()))
		return
	}
	
	var st reflect.Type
	var interfaceMethods map[string]bool
	if service.Interface == nil {
		logger.Warn("服务路由未绑定 Interface，将解析服务真实实现的导出方法未路由")
	} else {
		st = reflect.TypeOf(service.Interface)
		for {
			if st.Kind() != reflect.Ptr {
				break
			}
			st = st.Elem()
		}
		if st.NumMethod() == 0 {
			logger.Warn(fmt.Sprintf("接口 %s 没有可用的导出方法", st.String()))
			return
		}
		
		if !sv.Type().Implements(st) {
			err = fmt.Errorf("服务 %s 为实现接口 %s", sv.String(), st.String())
			return
		}
		interfaceMethods = map[string]bool{}
		for i := 0; i < st.NumMethod(); i++ {
			interfaceMethods[st.Method(i).Name] = true
		}
	}
	for i := 0; i < sv.NumMethod(); i++ {
		mt := sv.Type().Method(i)
		if !mt.IsExported() || !p.isInterfaceMethod(interfaceMethods, mt.Name) {
			continue
		}
		out = append(out, Route{
			Path:    fmt.Sprintf("/%s", sv.Type().Method(i).Name),
			Method:  Post,
			handler: sv.Method(i),
		})
	}
	
	out, err = p.prepareRoutes(out)
	if err != nil {
		return
	}
	
	return
}

func (p *App) isInterfaceMethod(interfaceMethods map[string]bool, name string) bool {
	if interfaceMethods == nil {
		return true
	}
	_, ok := interfaceMethods[name]
	return ok
}

// 递归注册路由树，处理中间件前缀逻辑，代理路由处理器为 Gin 控制器
func (p *App) registerRoutes(register Register, prefix string, routes []Route) (err error) {
	for _, v := range routes {
		if !v.handler.IsValid() {
			err = p.registerRoutes(
				register.Group(v.Prefix, v.Middlewares...),
				strings.Join([]string{prefix, v.Prefix}, ""),
				v.Routes,
			)
			if err != nil {
				return
			}
		} else {
			info := p.parseHandlerInfo(v.Handler)
			path := info.ParsePath(v.Path)
			p.addMethod(v.Method, strings.Join([]string{prefix, path}, ""), v.Description, info, v.handler)
			switch v.Method {
			case http.MethodGet:
				register.GET(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodPost:
				register.POST(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodPut:
				register.PUT(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodDelete:
				register.DELETE(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodHead:
				register.HEAD(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodOptions:
				register.OPTIONS(path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			default:
				err = fmt.Errorf("unsupport method: %s", v.Method)
				return
			}
			fmt.Println("处理器函数信息：", methodInfo(v.handler).Name())
			//p.handlers[v.handler.Type().String()]
		}
	}
	return
}

func (p *App) proxyHandler(handler reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		var out []reflect.Value
		ctx := context.Background()
		if handler.Type().NumIn() == 2 {
			var in reflect.Value
			var err error
			in, err = bind(c, handler.Type().In(1))
			if err != nil {
				_ = c.Error(err)
				return
			}
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx), in})
		} else {
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
		l := len(out)
		if !out[l-1].IsNil() {
			c.JSON(http.StatusInternalServerError, &Status{
				Detail: out[l-1].Interface().(error).Error(),
			})
			return
		}
		if l == 2 {
			switch out[0].Interface().(type) {
			case Html:
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(out[0].Interface().(Html)))
			case Text:
				c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(out[0].Interface().(Text)))
			default:
				c.JSON(http.StatusOK, out[0].Interface())
			}
			return
		} else {
			c.String(http.StatusOK, "")
			return
		}
	}
}

func isBasicType(v reflect.Value) bool {
	for {
		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func realType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t
		}
		t = t.Elem()
	}
}

func bind(c *gin.Context, t reflect.Type) (reflect.Value, error) {
	ptr := t.Kind() == reflect.Ptr
	if ptr {
		t = realType(t)
	}
	in := reflect.New(t)
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.MustBindWith(in.Interface(), b)
	if err != nil {
		return in, err
	}
	if ptr {
		return in, nil
	}
	return in.Elem(), nil
}

// 解析 Handler 的信息
func (p *App) parseHandlerInfo(h interface{}) HandlerInfo {
	target := reflect.ValueOf(h).Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

// 解析 Handler 的信息
func (p *App) parseHandlerInfoValue(v reflect.Value) HandlerInfo {
	target := v.Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

func (p *App) addMethod(method, path, description string, info HandlerInfo, handler reflect.Value) {
	if p.exporter == nil {
		return
	}
	m := &exporter.Method{
		Name:        info.Name,
		Path:        path,
		Method:      method,
		Description: description,
	}
	if handler.Type().NumIn() > 1 {
		m.Input = p.exporter.ReactField(
			handler.Type().In(1),
			"", "", "", nil)
	}
	if handler.Type().NumOut() > 1 {
		m.Output = p.exporter.ReactField(
			handler.Type().Out(0),
			"", "", "", nil)
	}
	p.methods = append(p.methods, m)
}
