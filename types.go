package goalfa

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

const (
	Post   = "POST"
	Get    = "GET"
	Delete = "DELETE"
	Put    = "PUT"
)

type Html string
type Text string

type Register interface {
	DELETE(path string, handles ...gin.HandlerFunc) gin.IRoutes
	GET(path string, handles ...gin.HandlerFunc) gin.IRoutes
	HEAD(path string, handles ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PATCH(path string, handles ...gin.HandlerFunc) gin.IRoutes
	POST(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PUT(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Any(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Group(prefix string, middleware ...gin.HandlerFunc) *gin.RouterGroup
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
}

type Param struct {
	Field       string
	Description string
}

type Service struct {
	Interface interface{}
	Instance  interface{}
}

type Route struct {
	Path        string
	Method      string
	Prefix      string
	Description string
	Middlewares []gin.HandlerFunc `json:"-"`
	Routes      []Route
	Params      []Param
	Service     Service     `json:"-"`
	Handler     interface{} `json:"-"`
	handler     reflect.Value
	handlerInfo HandlerInfo
}

func (r *Route) BindParams(params ...Param) {

}

type Router interface {
	Routes() []Route
}

type Module struct {
}

type Driver interface {
	Register(api *App) error
	Start(addr string) error
}

type HandlerInfo struct {
	Name     string
	Location string
}

func (p HandlerInfo) ParsePath(path string) string {
	if path != "" {
		return path
	}
	return fmt.Sprintf("/%s", p.Name)
}

type RouteTable struct {
	rows []RouteRow
}

func (p *RouteTable) AddRow(method, path string, info HandlerInfo) {
	p.rows = append(p.rows, RouteRow{
		Method:   method,
		Path:     path,
		Source:   info.Name,
		Location: info.Location,
	})
}

func (p RouteTable) Print() {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"", "接口地址", "方法名", "位置"})
	var data [][]string
	for _, v := range p.rows {
		data = append(data, []string{v.Method, v.Path, v.Source, v.Location})
	}
	t.AppendBulk(data)
	t.Render()
}

type RouteRow struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Source   string `json:"source"`
	Location string `json:"location"`
}

type Value struct {
	value   interface{}
	headers map[string]string
}
