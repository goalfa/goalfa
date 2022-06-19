package exporter

import "fmt"

type Options struct {
	Project    string           `json:"project"`
	Envs       []Env            `json:"envs"`
	BasicTypes []BasicType      `json:"-"`
	Makers     map[string]Maker `json:"-"`
}

type BasicType struct {
	Elem    interface{}        `json:"-"`
	Type    string             `json:"type"`
	Package string             `json:"package,omitempty"`
	Mapping map[string]Library `json:"structsMapping,omitempty"`
}

func (p BasicType) Fork() *BasicType {
	n := new(BasicType)
	n.Elem = p.Elem
	if p.Mapping != nil {
		n.Mapping = map[string]Library{}
		for k, v := range p.Mapping {
			n.Mapping[k] = v
		}
	}
	n.Package = p.Package
	return n
}

func (p BasicType) getMapping(lang string) *Library {
	if p.Mapping == nil {
		return nil
	}
	v, ok := p.Mapping[lang]
	if !ok {
		return nil
	}
	return &v
}

type BasicTypes struct {
	list    []*BasicType
	mapping map[string]bool
}

func (p *BasicTypes) Add(item *BasicType) {
	if p.mapping == nil {
		p.mapping = map[string]bool{}
	}
	if _, ok := p.mapping[item.Type]; ok {
		return
	}
	p.mapping[item.Type] = true
	p.list = append(p.list, item)
}

func (p BasicTypes) All() []*BasicType {
	return p.list
}

type Library struct {
	Type    string   `json:"type,omitempty"`
	Package *Package `json:"package,omitempty"`
}

type Package struct {
	Import string `json:"import,omitempty"`
	From   string `json:"from,omitempty"`
}

type Env struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type Method struct {
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Middlewares string `json:"middlewares,omitempty"`
	Input       *Field `json:"input,omitempty"`
	Output      *Field `json:"output,omitempty"`
}

func (p Method) Fork() *Method {
	n := new(Method)
	n.Name = p.Name
	n.Path = p.Path
	n.Method = p.Method
	n.Description = p.Description
	n.Middlewares = p.Middlewares
	if p.Input != nil {
		n.Input = p.Input.Fork()
	}
	if p.Output != nil {
		n.Output = p.Output.Fork()
	}
	return n
}

type Struct struct {
	Name    string   `json:"name"`
	Package string   `json:"package"`
	Fields  []*Field `json:"fields"`
}

type Field struct {
	Name        string `json:"name,omitempty"`
	Param       string `json:"param,omitempty"`
	Label       string `json:"label,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	IsArray     bool   `json:"array,omitempty"`
	IsStruct    bool   `json:"struct,omitempty"`
	IsBasic     bool   `json:"isBasic"`
	Origin      string `json:"origin,omitempty"` // 原始类型
	//Fields      []*Field   `json:"fields,omitempty"`    // 描述 IsStruct 成员变量
	Elem      *Field     `json:"elem,omitempty"`      // 描述 Slice/IsArray 子元素
	Validator *Validator `json:"validator,omitempty"` // 定义校验器
	Form      string     `json:"form,omitempty"`      // 定义表单组件
}

func (p Field) Fork() *Field {
	n := new(Field)
	n.Name = p.Name
	n.Param = p.Param
	n.Label = p.Label
	n.Type = p.Type
	n.Description = p.Description
	n.IsArray = p.IsArray
	n.IsStruct = p.IsStruct
	n.IsBasic = p.IsBasic
	n.Origin = p.Origin
	//for _, v := range p.Fields {
	//	n.Fields = append(n.Fields, v.Fork())
	//}
	if p.Elem != nil {
		n.Elem = p.Elem.Fork()
	}
	n.Validator = p.Validator
	n.Form = p.Form

	return n
}

type Structs struct {
	structs        []*Struct
	structsMapping map[string]*Struct
	namesMapping   map[string]int
}

func (p *Structs) GetStruct(pkg, name string) *Struct {
	key := p.Key(pkg, name)
	if v, ok := p.structsMapping[key]; ok {
		return v
	}
	return nil
}

func (p *Structs) GetStructName(name string) string {
	if p.namesMapping == nil {
		return name
	}
	count, ok := p.namesMapping[name]
	if !ok || count == 0 {
		return name
	}
	return fmt.Sprintf("%s%d", name, count)
}

func (p *Structs) Key(pkg, name string) string {
	return fmt.Sprintf("%s@%s", pkg, name)
}

func (p *Structs) Add(item *Struct) {
	if p.structsMapping == nil {
		p.structsMapping = map[string]*Struct{}
	}
	if p.namesMapping == nil {
		p.namesMapping = map[string]int{}
	}
	key := p.Key(item.Name, item.Package)
	if _, ok := p.structsMapping[key]; ok {
		return
	}
	p.structsMapping[key] = item
	count, ok := p.namesMapping[item.Name]
	if !ok {
		p.namesMapping[item.Name] = 0
		count = 0
	}
	p.namesMapping[item.Name] = count + 1
	p.structs = append(p.structs, item)
}

func (p Structs) All() []*Struct {
	return p.structs
}

type Validator struct {
	Required bool     `json:"required,omitempty"`
	Max      *uint64  `json:"max,omitempty"`
	Min      *int64   `json:"min,omitempty"`
	Enums    []string `json:"enums,omitempty"`
}

type Component struct {
	Name string
}
