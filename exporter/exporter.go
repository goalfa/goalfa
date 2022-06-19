package exporter

import (
	"fmt"
	"github.com/datafony/alfa/assets"
	"github.com/datafony/alfa/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_log"
	"github.com/gozelle/_log/wrap"
	"github.com/ttacon/chalk"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func NewExporter(addr string, options *Options) *Exporter {
	e := &Exporter{addr: addr, options: options}
	e.structs = new(Structs)
	e.initBasicTypes()
	e.initMakers()
	return e
}

type Exporter struct {
	version string
	addr    string
	options *Options
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

func (p Exporter) Run() {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
	Version string       `json:"version"`
	Options *Options     `json:"options"`
	Methods []*Method    `json:"methods"`
	Basics  []*BasicType `json:"basics,omitempty"`
	Structs []*Struct    `json:"structs"`
}

// 导出接口描述协议
func (p Exporter) protocolHandler(c *gin.Context) {
	out := new(ProtocolOutput)
	out.Version = p.version
	out.Options = p.options
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
	//for _, v := range field.Fields {
	//	p.toTsProtocolFieldType(v)
	//}
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

// ReflectFields 反射转换输入输出的字段信息
//func (p *Exporter) ReflectFields(name, param, label string, validator *Validator, t reflect.Name) (field *Field) {
//
//	t = utils.TypeElem(t)
//
//	typ := p.reflectFieldType(t)
//	pkg := t.PkgPath()
//
//	if p.structs.HasStruct(typ, pkg) {
//		return
//	}
//
//	field = new(Field)
//	field.Name = typ
//	field.Name = name
//	field.Param = param
//	field.Label = label
//	//field.Package = pkg
//	field.Validator = validator
//
//	//p.structs.Add(field)
//
//	if t.Kind() == reflect.IsStruct && p.getBasicType(t) == nil {
//		field.IsStruct = true
//		for i := 0; i < t.NumField(); i++ {
//			sf := t.Field(i)
//			_name := sf.Name
//			_param := p.getParam(sf)
//			_label := p.getFieldLabel(sf)
//			_validator := p.getFieldValidator(sf)
//			_field := p.ReflectFields(_name, _param, _label, _validator, sf.Name)
//			if _field != nil && (_field.IsStruct || _field.Nested) {
//				field.Nested = true
//			}
//			field.Fields = append(field.Fields, _field)
//		}
//	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.IsArray {
//		field.IsArray = true
//		field.Elem = p.ReflectFields("", "", label, validator, t.Elem())
//		if field.Elem != nil && (field.Elem.IsStruct || field.Elem.Nested) {
//			field.Nested = true
//		}
//	}
//
//	return
//}

func (p *Exporter) ReactField(elem reflect.Type, name, param, label string, validator *Validator) *Field {

	elem = utils.TypeElem(elem)
	fieldType, isBasic, isSturct := p.reflectFieldType(elem)
	//fieldPackage := elem.PkgPath()

	//if p.structs.HasStruct(typ, pkg) {
	//	return
	//}

	field := new(Field)
	field.Name = name
	field.Param = param
	field.Label = label
	//field.Package = fieldPackage
	field.Validator = validator
	field.Type = fieldType
	field.IsStruct = isSturct
	field.IsBasic = isBasic

	return field
}

func (p *Exporter) reflectFieldType(elem reflect.Type) (typeName string, isBasic, isStruct bool) {
	basicType := p.getBasicType(elem)
	if basicType != nil {
		return elem.String(), true, false
	}
	if elem.Kind() == reflect.Struct {
		return p.ReflectStruct(elem, 0), false, true
	}
	return p.getType(elem), false, false
}

// ReflectStruct 反射转换输入输出的字段信息
func (p *Exporter) ReflectStruct(elem reflect.Type, depth int) string {

	elem = utils.TypeElem(elem)

	if elem.Kind() != reflect.Struct {
		return ""
	}

	name := elem.Name()
	pkg := elem.PkgPath()

	if v := p.structs.GetStruct(pkg, name); v != nil {
		return v.Name
	}

	s := new(Struct)
	s.Name = p.structs.GetStructName(name)
	s.Package = pkg
	p.structs.Add(s)
	for i := 0; i < elem.NumField(); i++ {
		item := elem.Field(i)
		param := p.getParam(item)
		label := p.getFieldLabel(item)
		validator := p.getFieldValidator(item)
		s.Fields = append(
			s.Fields,
			p.ReactField(item.Type, item.Name, param, label, validator),
		)
	}
	return s.Name
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
