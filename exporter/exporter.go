package exporter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_log"
	"github.com/gozelle/_log/wrap"
	"github.com/koyeo/goalfa/assets"
	"github.com/koyeo/goalfa/utils"
	"github.com/ttacon/chalk"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func NewExporter(addr string, options *Settings) *Exporter {
	e := &Exporter{addr: addr, options: options}
	e.structs = new(Structs)
	e.initBasicTypes()
	e.initMakers()
	return e
}

type Exporter struct {
	version string
	addr    string
	options *Settings
	Name    string
	Package string
	methods []*Method
	basics  map[string]*BasicType
	structs *Structs
	makers  map[string]Maker
}

func (p *Exporter) Init(version string, methods []*Method, models *Structs) {
	p.version = version
	p.methods = methods
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		
		c.Next()
	}
}

func (p Exporter) Run() {
	// 自定义 gin 驱动
	engine := gin.Default()
	//engine.Use(cors.New(config))
	engine.Use(Cors())
	//engine.
	gin.SetMode(gin.DebugMode)
	engine.GET("/sdk", p.sdkHandler)
	engine.GET("/protocol", p.protocolHandler)
	engine.StaticFS("/exporter", assets.Root)
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/exporter/index.html")
	})
	engine.GET("/:path", func(c *gin.Context) {
		path := c.Param("path")
		if path == "exporter" {
			path = "index.html"
		}
		if !strings.HasPrefix(path, "exporter") {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/exporter/%s", path))
		}
	})
	engine.GET("/test", func(c *gin.Context) {
		c.Request.URL.Path = "/sdk"
		engine.HandleContext(c)
	})
	go func() {
		p.printAddress()
		err := engine.Run(p.addr)
		if err != nil {
			_log.Panic("接口导出器启动失败", wrap.Error(err))
		}
	}()
}

// 打印 API 调试器访问地址
func (p Exporter) printAddress() {
	addr := p.addr
	if strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf("http://127.0.0.1%s", addr)
	} else {
		addr = fmt.Sprintf("http://%s", addr)
	}
	fmt.Println(chalk.Green.Color(strings.Repeat("=", 100)))
	fmt.Println(chalk.Green.Color(fmt.Sprintf("API 调试器访问地址：%s", addr)))
	fmt.Println(chalk.Green.Color(strings.Repeat("=", 100)))
}

// 导出 SDK 代码
func (p Exporter) sdkHandler(c *gin.Context) {
	sdk := NewSDK(p.methods)
	data, err := sdk.Make(p.makers, c.Query("lang"), c.Query("package"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

type ProtocolOutput struct {
	Version  string       `json:"version"`
	Settings *Settings    `json:"options"`
	Methods  []*Method    `json:"methods"`
	Basics   []*BasicType `json:"basics,omitempty"`
	Structs  []*Struct    `json:"structs"`
}

// 导出接口描述协议
func (p Exporter) protocolHandler(c *gin.Context) {
	out := new(ProtocolOutput)
	out.Version = p.version
	out.Settings = p.options
	out.Methods = p.convertMethodTypes(c.Query("lang"))
	basics := new(BasicTypes)
	for _, v := range p.basics {
		basics.Add(v)
	}
	out.Basics = basics.All()
	out.Structs = p.structs.All()
	c.JSON(200, out)
}

func (p Exporter) convertMethodTypes(lang string) []*Method {
	methods := make([]*Method, 0)
	switch lang {
	case "ts":
		for _, v := range p.methods {
			n := v.Fork()
			p.toTsProtocolFieldType(n.Input)
			p.toTsProtocolFieldType(n.Output)
			methods = append(methods, n)
		}
	default:
		methods = p.methods
	}
	return methods
}

func (p Exporter) toTsProtocolFieldType(field *Field) {
	if field == nil {
		return
	}
	field.Origin = field.Type
	field.Type = tsProtocolTypeConverter(field.Type)
	for _, v := range field.Fields {
		p.toTsProtocolFieldType(v)
	}
	if field.Elem != nil {
		p.toTsProtocolFieldType(field.Elem)
	}
}

func (p *Exporter) initBasicTypes() {
	if p.options == nil {
		return
	}
	basicTypes := []BasicType{}
	basicTypes = append(basicTypes, p.options.BasicTypes...)
	for _, v := range basicTypes {
		if p.basics == nil {
			p.basics = map[string]*BasicType{}
		}
		r := reflect.ValueOf(v.Elem)
		basicType := v.Fork()
		basicType.Package = r.Type().PkgPath()
		basicType.Type = r.Type().String()
		p.basics[fmt.Sprintf("%s@%s", basicType.Package, r.Type().String())] = basicType
	}
}

func (p *Exporter) initMakers() {
	p.makers = map[string]Maker{
		"go":      GoMaker{},
		"angular": AngularMaker{},
	}
	if p.options != nil {
		for k, v := range p.options.Makers {
			p.makers[k] = v
		}
	}
}

func (p *Exporter) ReactField(elem reflect.Type, name, param, label string, validator *Validator) *Field {
	
	elem = utils.TypeElem(elem)
	
	field := new(Field)
	field.Name = name
	field.Param = param
	field.Label = label
	field.Validator = validator
	fieldType, isBasic, isStruct, isArray := p.reflectFieldType(field, elem)
	field.Type = fieldType
	field.IsStruct = isStruct
	field.IsBasic = isBasic
	field.IsArray = isArray
	
	return field
}

func (p *Exporter) reflectFieldType(field *Field, elem reflect.Type) (typeName string, isBasic, isStruct, isArray bool) {
	basicType := p.getBasicType(elem)
	if basicType != nil {
		return elem.String(), true, false, false
	} else if elem.Kind() == reflect.Struct {
		return p.ReflectStruct(field, elem, 0), false, true, false
	} else if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
		_elem := utils.TypeElem(elem.Elem())
		_field := p.ReactField(_elem, "", "", "", nil)
		return _field.Type, false, true, true
	}
	return p.getType(elem), false, false, false
}

// ReflectStruct 反射转换输入输出的字段信息
func (p *Exporter) ReflectStruct(field *Field, elem reflect.Type, depth int) string {
	
	elem = utils.TypeElem(elem)
	
	name := elem.Name()
	pkg := elem.PkgPath()
	fields := make([]*Field, 0)
	if name != "" {
		if v := p.structs.GetStruct(pkg, name); v != nil {
			return v.Name
		}
	}
	for i := 0; i < elem.NumField(); i++ {
		item := elem.Field(i)
		param := p.getParam(item)
		label := p.getFieldLabel(item)
		validator := p.getFieldValidator(item)
		fields = append(
			fields,
			p.ReactField(item.Type, item.Name, param, label, validator),
		)
	}
	if name != "" {
		s := new(Struct)
		s.Name = p.structs.GetStructName(name)
		s.Package = pkg
		p.structs.Add(s)
		s.Fields = fields
		return s.Name
	}
	field.Fields = fields
	return ""
}

func (p Exporter) getBasicType(t reflect.Type) *BasicType {
	if p.basics == nil {
		return nil
	}
	v, ok := p.basics[fmt.Sprintf("%s@%s", t.PkgPath(), t.String())]
	if !ok {
		return nil
	}
	return v
}

func (p Exporter) getType(t reflect.Type) string {
	s := t.String()
	if strings.Contains(s, ".") {
		s = strings.Split(s, ".")[1]
	}
	return s
}

func (p Exporter) getFieldLabel(field reflect.StructField) string {
	return field.Tag.Get("label")
}

func (p Exporter) getFieldValidator(field reflect.StructField) (validator *Validator) {
	required := strings.Contains(field.Tag.Get("validator"), "required")
	if required {
		validator = p.newIfNoValidator(validator)
		validator.Required = true
	}
	return
}

func (p Exporter) newIfNoValidator(validator *Validator) *Validator {
	if validator == nil {
		validator = new(Validator)
	}
	return validator
}

func (p Exporter) getParam(field reflect.StructField) string {
	n := field.Tag.Get("json")
	return strings.ReplaceAll(n, ",omitempty", "")
}
