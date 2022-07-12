package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goalfa "github.com/koyeo/goalfa"
	"github.com/koyeo/goalfa/example/foo/foo-service"
	"github.com/koyeo/goalfa/exporter"
)

func main() {
	
	// gin 跨域配置
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowAllOrigins = true
	
	// 自定义 gin 驱动
	engine := gin.Default()
	engine.Use(cors.New(config))
	
	// goalfa 实例
	app := goalfa.New()
	app.SetVersion("1.0.0")
	app.SetEngine(engine)
	app.AddRouter(foo_service.NewFooRouter(new(foo_service.FooMock)))
	
	// App 导出器配置
	app.SetExporter(":9090", &exporter.Settings{
		Project: "Foo",
		Envs: []exporter.Env{
			{
				Name: "本地测试",
				Host: "http://localhost:8080",
			},
		},
		Makers: map[string]exporter.Maker{
			"python": exporter.GoMaker{},
		},
		//BasicTypes: []exporter.BasicType{
		//	{
		//		Elem: service.CID{},
		//		Mapping: map[string]exporter.Library{
		//			"ts": {Type: "string"},
		//		},
		//	},
		//},
	})
	
	app.Run(":8080")
}
