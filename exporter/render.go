package exporter

import (
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v5"
	"strings"
)

type RenderData struct {
	Packages []*RenderPackage
	Methods  []*RenderMethod
	Structs  []*RenderStruct
}

type RenderMethod struct {
	Name         string
	InputType    string
	OutputType   string
	OutputStruct bool
	Description  string
	Method       string
	Path         string
}

type RenderStruct struct {
	Name        string
	Type        string
	Description string
	Fields      []*RenderField
}

type RenderField struct {
	Name        string
	Param       string
	Type        string
	Description string
	Required    bool
	Label       string
}

type RenderPackage struct {
	Import        string
	From          string
	importList    []string
	importMapping map[string]bool
}

func (p *RenderPackage) addImport(i string) {
	if p.importMapping == nil {
		p.importMapping = map[string]bool{}
	}
	_, ok := p.importMapping[i]
	if ok {
		return
	}
	p.importMapping[i] = true
	p.importList = append(p.importList, i)
}

type RenderPackages struct {
	list    []*RenderPackage
	mapping map[string]*RenderPackage
}

func (p *RenderPackages) add(item *RenderPackage) {
	if p.mapping == nil {
		p.mapping = map[string]*RenderPackage{}
	}
	_, ok := p.mapping[item.From]
	if ok {
		p.mapping[item.From].addImport(item.Import)
		return
	}
	p.mapping[item.From] = item
	p.list = append(p.list, item)
}

type Namer func(string) string
type Formatter func(string) (string, error)
type Typer func(s string, isStruct, isArray bool) string

var EmptyNamer Namer = func(s string) string {
	return s
}

var EmptyTyper Typer = func(s string, isStruct, isArray bool) string {
	return s
}

var EmptyFormatter Formatter = func(s string) (string, error) {
	return s, nil
}

func MakeRenderData(lang string, methods []*Method, namer Namer, typer Typer) (data *RenderData) {
	data = new(RenderData)
	checker := newRenderFieldChecker()
	renderPackages := new(RenderPackages)
	for _, v := range methods {
		data.Methods = append(data.Methods, makeRenderMethod(lang, v, namer, typer, renderPackages))
		data.Structs = append(data.Structs, makeRenderStructs(lang, v, namer, typer, checker, renderPackages)...)
		data.Packages = renderPackages.list
	}
	return
}

func makeRenderMethod(lang string, method *Method, namer Namer, typer Typer, renderPackages *RenderPackages) (renderMethod *RenderMethod) {
	renderMethod = new(RenderMethod)
	renderMethod.Name = namer(method.Name)
	renderMethod.Description = method.Description
	renderMethod.Method = method.Method
	renderMethod.Path = method.Path
	if method.Input != nil {
		renderMethod.InputType = makeMethodIOName(lang, method.Input, typer, renderPackages)
	}
	if method.Output != nil {
		renderMethod.OutputType = makeMethodIOName(lang, method.Output, typer, renderPackages)
		renderMethod.OutputStruct = method.Output.IsStruct
	}
	return
}

func makeMethodIOName(lang string, field *Field, typer Typer, renderPackages *RenderPackages) string {
	if field.IsArray {
		return parseNestedType(lang, field, typer, renderPackages)
	} else if field.IsStruct {
		return typer(getRenderFieldType(lang, field, renderPackages), true, false)
	} else {
		return typer(getRenderFieldType(lang, field, renderPackages), false, false)
	}
}

func parseNestedType(lang string, field *Field, typer Typer, renderPackages *RenderPackages) string {
	if field.IsArray {
		return typer(parseNestedType(lang, field.Elem, typer, renderPackages), field.IsStruct, field.IsArray)
	} else {
		return typer(getRenderFieldType(lang, field, renderPackages), field.IsStruct, field.IsArray)
	}
}

func getRenderFieldType(lang string, field *Field, renderPackages *RenderPackages) string {
	_type := field.Type
	//if field.IsBasic != nil {
	//	if lang == Go {
	//		renderPackages.add(&RenderPackage{
	//			From: field.IsBasic.Package,
	//		})
	//	} else {
	//		lib := field.IsBasic.getMapping(lang)
	//		if lib != nil {
	//			_type = lib.Name
	//			if lib.Package != nil && lib.Package.From != "" {
	//				if renderPackages != nil {
	//					renderPackages.add(&RenderPackage{
	//						Import: lib.Package.Import,
	//						From:   lib.Package.From,
	//					})
	//				}
	//			}
	//		}
	//	}
	//}
	return _type
}

func makeRenderStructs(lang string, method *Method, namer Namer, typer Typer,
	checker *renderFieldChecker, renderPackages *RenderPackages) (renderFields []*RenderStruct) {
	//if method.Input != nil && (method.Input.IsStruct || (method.Input.IsArray && method.Input.Nested)) {
	//	toRenderStructs(lang, method.Input, namer, typer, checker, &renderFields, renderPackages)
	//}
	//if method.Output != nil && (method.Output.IsStruct || (method.Output.IsArray && method.Output.Nested)) {
	//	toRenderStructs(lang, method.Output, namer, typer, checker, &renderFields, renderPackages)
	//}
	return
}

func toRenderStructs(lang string, field *Field, namer Namer, typer Typer, checker *renderFieldChecker,
	renderStructs *[]*RenderStruct, renderPackages *RenderPackages) {
	//if field.IsArray && field.Nested { // 处理嵌套数组对象
	//	toRenderStructs(lang, field.Elem, namer, typer, checker, renderStructs, renderPackages)
	//} else if field.IsStruct { // 处理对象
	//	renderStruct := new(RenderStruct)
	//	name := namer(field.Name)
	//	renderStruct.Name = name
	//	renderStruct.Description = field.Description
	//	if !checker.HasStruct(name) {
	//		checker.Add(name)
	//		*renderStructs = append(*renderStructs, renderStruct)
	//	}
	//	for _, v := range field.Fields {
	//		renderField := new(RenderField)
	//		renderField.Name = namer(v.Name)
	//		renderField.Param = v.Param
	//		renderField.Name = parseNestedType(lang, v, typer, renderPackages)
	//		renderField.Description = v.Description
	//		renderField.Label = v.Label
	//		if v.Validator != nil {
	//			renderField.Required = v.Validator.Required
	//		}
	//		renderStruct.Fields = append(renderStruct.Fields, renderField)
	//		if v.IsStruct {
	//			toRenderStructs(lang, v, namer, typer, checker, renderStructs, renderPackages)
	//		} else if v.IsArray && v.Nested {
	//			toRenderStructs(lang, v.Elem, namer, typer, checker, renderStructs, renderPackages)
	//		}
	//	}
	//}
	return
}

func Render(tpl string, data interface{}, formatter Formatter) (result string, err error) {
	_tpl, err := pongo2.FromString(tpl)
	if err != nil {
		return
	}
	ctx := structs.Map(data)
	ctx["_trimPrefix"] = strings.TrimPrefix
	result, err = _tpl.Execute(ctx)
	if err != nil {
		return
	}
	result, err = formatter(result)
	if err != nil {
		return
	}
	return
}

func newRenderFieldChecker() *renderFieldChecker {
	return &renderFieldChecker{
		cache: map[string]bool{},
	}
}

type renderFieldChecker struct {
	cache map[string]bool
}

func (p *renderFieldChecker) Has(name string) bool {
	_, ok := p.cache[name]
	return ok
}

func (p *renderFieldChecker) Add(name string) {
	p.cache[name] = true
}
