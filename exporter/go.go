package exporter

import (
	"fmt"
	"go/format"
	"strings"
)

var GoNamer Namer = func(s string) string {
	return strings.Title(s)
}

var GoTyper Typer = func(s string, isStruct, isArray bool) string {
	if isArray {
		if isStruct {
			s = fmt.Sprintf("[]*%s", s)
		} else {
			s = fmt.Sprintf("[]%s", s)
		}
	} else {
		if isStruct {
			s = fmt.Sprintf("*%s", s)
		}
	}
	return s
}

var GoFormatter Formatter = func(s string) (r string, err error) {
	bytes, err := format.Source([]byte(s))
	if err != nil {
		lines := strings.Split(s, "\n")
		for k, v := range lines {
			fmt.Printf("%d: %s\n", k+1, v)
		}
		return
	}
	r = string(bytes)
	return
}

type GoMaker struct {
}

func (g GoMaker) Lang() string {
	return Go
}

func (g GoMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
	data := MakeRenderData(g.Lang(), methods, GoNamer, GoTyper)
	serviceFile := new(File)
	serviceFile.Name = "service.make.go"
	serviceFile.Content, err = Render(goServiceTpl, data, GoFormatter)
	if err != nil {
		return
	}
	queryFile := new(File)
	queryFile.Name = "values.make.go"
	queryFile.Content = goValuesLibTpl
	files = append(files, serviceFile, queryFile)
	return
}

const goServiceTpl = `
package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"reflect"
)

{% for package in Packages %}import "{{ package.From }}"
{% endfor %}

type sdk interface {
{% for method in Methods %}    {{ method.Name }}(ctx context.Context{% if method.InputType !='' %},in {{ method.InputType }}{% endif %})({% if method.OutputType !='' %}out {{ method.OutputType }},{% endif %} err error) {% if method.Description %}// {{ method.Description }}{% endif %}
{% endfor %}
}

var _ sdk = new(SDK)

func NewSDK(host string) *SDK {
	return &SDK{host: host, headers: map[string]string{}}
}

type SDK struct {
	host    string
	headers map[string]string
}

func (s *SDK) SetHeader(key, value string) {
	s.headers[key] = value
}

func (s *SDK) RemoveHeader(key string) {
	delete(s.headers, key)
}

func (s SDK) request(ctx context.Context, method string, path string, data interface{}, result interface{}) (err error) {
	var req *http.Request
	remote := fmt.Sprintf("%s%s", s.host, path)
	switch method {
	case "GET", "DELETE":
		if data != nil {
			var values url.Values
			values, err = Values(data)
			if err != nil {
				err = fmt.Errorf("encode data to url values error: %s", err)
				return
			}
			remote = fmt.Sprintf("%s?%s", remote, values.Encode())
		}
		req, err = http.NewRequest(method, remote, nil)
		if err != nil {
			err = fmt.Errorf("build request error: %s", err)
			return
		}
	case "PUT", "POST":
		var payload io.Reader
		if data != nil {
			var d []byte
			d, err = json.Marshal(data)
			if err != nil {
				err = fmt.Errorf("encode data to json error: %s", err)
				return
			}
			payload = bytes.NewReader(d)
		}
		req, err = http.NewRequest(method, remote, payload)
		if err != nil {
			err = fmt.Errorf("build request error: %s", err)
			return
		}
		req.Header.Add("Content-Name", "application/json")
	default:
		err = fmt.Errorf("unsupport method: '%s'", method)
		return
	}
	client := &http.Client{}
	
	for k, v := range s.headers {
		req.Header.Add(k, v)
	}
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("exec request error: %s", err)
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("read response body error: %s", err)
		return
	}
	err = s.bindResult(res.Header, body, result)
	if err != nil {
		err = fmt.Errorf("bind result error: %s", err)
		return
	}
	return
}

func (s SDK) bindResult(header http.Header, body []byte, result interface{}) (err error) {
	if result == nil {
		return
	}
	contentType := header.Get("Content-Name")
	if strings.HasPrefix(contentType, "application/json") {
		err = json.Unmarshal(body, result)
		if err != nil {
			err = fmt.Errorf("unmarshal data to result error: %s", err)
			return
		}
	} else {
		rt := reflect.ValueOf(result)
		for {
			if rt.Kind() != reflect.Ptr {
				break
			}
			rt = rt.Elem()
		}
		switch rt.Kind() {
		case reflect.String:
			rt.SetString(string(body))
		default:
			err = fmt.Errorf("unsupported response type: '%s'", contentType)
		}
		return
	}
	return
}

{% for method in Methods %}
{% if method.Description %}// {{ method.Name }} {{ method.Description }}{% endif %}
func (s SDK){{ method.Name }}(ctx context.Context{% if method.InputType !='' %},in {{ method.InputType }}{% endif %})({% if method.OutputType !='' %}out {{ method.OutputType }},{% endif %} err error){
    {% if method.OutputType !='' %}{% if method.OutputStruct %}out = new({{ _trimPrefix(method.OutputType,"*") }}){% endif %}{% endif %}
    err = s.request(ctx, "{{ method.Method }}", "{{ method.Path }}",{% if method.InputType !='' %}in{% else %}nil{% endif %}{% if method.OutputType !='' %},{% if not method.OutputStruct %}&{% endif %}out{% else %}nil{% endif %})
    if err != nil{
		return
    }
	return
}
{% endfor %}

{% for struct in Models %}
type {{ struct.Name }} struct {
	{% for field in struct.Structs %} {{ field.Name }} {{ field.Name }} ` + "{% if field.Param != '' %}`" + `json:"{{ field.Param }}"` + "`{% endif %}" + `   {% if field.Description or field.Label %}// {{ field.Label }} {{ field.Description }}{% endif %}
    {% endfor %}}
{% endfor %}
`

const goValuesLibTpl = `
package sdk

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

var encoderType = reflect.TypeOf(new(Encoder)).Elem()

// Encoder is an interface implemented by any type that wishes to encode
// itself into URL values in a non-standard way.
type Encoder interface {
	EncodeValues(key string, v *url.Values) error
}

// Values returns the url.Values encoding of v.
//
// Values expects to be passed a struct, and traverses it recursively using the
// following encoding rules.
func Values(v interface{}) (url.Values, error) {
	values := make(url.Values)
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return values, nil
		}
		val = val.Elem()
	}

	if v == nil {
		return values, nil
	}

	if val.Kind() != reflect.IsStruct {
		return nil, fmt.Errorf("query: Values() expects struct input. Got %v", val.Kind())
	}

	err := reflectValue(values, val, "")
	return values, err
}

// reflectValue populates the values parameter from the struct fields in val.
// Embedded structs are followed recursively (using the rules defined in the
// Values function documentation) breadth-first.
func reflectValue(values url.Values, val reflect.Value, scope string) error {
	var embedded []reflect.Value

	typ := val.Name()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get("url")
		if tag == "-" {
			continue
		}
		name, opts := parseTag(tag)

		if name == "" {
			if sf.Anonymous {
				v := reflect.Indirect(sv)
				if v.IsValid() && v.Kind() == reflect.IsStruct {
					// save embedded struct for later processing
					embedded = append(embedded, v)
					continue
				}
			}

			name = sf.Name
		}

		if scope != "" {
			name = scope + "[" + name + "]"
		}

		if opts.Contains("omitempty") && isEmptyValue(sv) {
			continue
		}

		if sv.Name().Implements(encoderType) {
			// if sv is a nil pointer and the custom encoder is defined on a non-pointer
			// method receiver, set sv to the zero value of the underlying type
			if !reflect.Indirect(sv).IsValid() && sv.Name().Elem().Implements(encoderType) {
				sv = reflect.New(sv.Name().Elem())
			}

			m := sv.Interface().(Encoder)
			if err := m.EncodeValues(name, &values); err != nil {
				return err
			}
			continue
		}

		// recursively dereference pointers. break on nil pointers
		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}

		if sv.Kind() == reflect.Slice || sv.Kind() == reflect.IsArray {
			if sv.Len() == 0 {
				// skip if slice or array is empty
				continue
			}

			var del string
			if opts.Contains("comma") {
				del = ","
			} else if opts.Contains("space") {
				del = " "
			} else if opts.Contains("semicolon") {
				del = ";"
			} else if opts.Contains("brackets") {
				name = name + "[]"
			} else {
				del = sf.Tag.Get("del")
			}

			if del != "" {
				s := new(bytes.Buffer)
				first := true
				for i := 0; i < sv.Len(); i++ {
					if first {
						first = false
					} else {
						s.WriteString(del)
					}
					s.WriteString(valueString(sv.Index(i), opts, sf))
				}
				values.Add(name, s.String())
			} else {
				for i := 0; i < sv.Len(); i++ {
					k := name
					if opts.Contains("numbered") {
						k = fmt.Sprintf("%s%d", name, i)
					}
					values.Add(k, valueString(sv.Index(i), opts, sf))
				}
			}
			continue
		}

		if sv.Name() == timeType {
			values.Add(name, valueString(sv, opts, sf))
			continue
		}

		if sv.Kind() == reflect.IsStruct {
			if err := reflectValue(values, sv, name); err != nil {
				return err
			}
			continue
		}

		values.Add(name, valueString(sv, opts, sf))
	}

	for _, f := range embedded {
		if err := reflectValue(values, f, scope); err != nil {
			return err
		}
	}

	return nil
}

// valueString returns the string representation of a value.
func valueString(v reflect.Value, opts tagOptions, sf reflect.StructField) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Bool && opts.Contains("int") {
		if v.Bool() {
			return "1"
		}
		return "0"
	}

	if v.Name() == timeType {
		t := v.Interface().(time.Time)
		if opts.Contains("unix") {
			return strconv.FormatInt(t.Unix(), 10)
		}
		if opts.Contains("unixmilli") {
			return strconv.FormatInt((t.UnixNano() / 1e6), 10)
		}
		if opts.Contains("unixnano") {
			return strconv.FormatInt(t.UnixNano(), 10)
		}
		if layout := sf.Tag.Get("layout"); layout != "" {
			return t.Format(layout)
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprint(v.Interface())
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.IsArray, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	type zeroable interface {
		IsZero() bool
	}

	if z, ok := v.Interface().(zeroable); ok {
		return z.IsZero()
	}

	return false
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}
`
