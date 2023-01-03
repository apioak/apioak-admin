package main

import (
	"apioak-admin/cores"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net"
)

//go:embed html/*
var htmlFS embed.FS

func main() {

	// 全局配置
	var conf cores.ConfigGlobal

	// 全局配置初始化
	if err := cores.InitConfig(&conf); err != nil {
		panic(err)
	}

	// 初始化框架
	if err := cores.InitFramework(&conf); err != nil {
		panic(err)
	}

	// 初始化Logger
	if err := cores.InitLogger(&conf); err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := cores.InitDataBase(&conf); err != nil {
		panic(err)
	}

	//初始化Token
	if err := cores.InitToken(&conf); err != nil {
		panic(err)
	}

	// 初始化参数验证器
	if err := cores.InitValidator(&conf); err != nil {
		panic(err)
	}

	// 初始化路由
	if err := cores.InitRouter(&conf); err != nil {
		panic(err)
	}

	initStatic(&conf)

	// 协程处理额外事件
	cores.InitGoroutineFunc()

	// 服务启动
	if err := cores.RunFramework(&conf); err != nil {
		panic(err)
	}
}

func initStatic(conf *cores.ConfigGlobal) {

	// 引入html
	conf.Runtime.Gin.SetHTMLTemplate(template.Must(template.New("").ParseFS(htmlFS, "html/*")))

	// 访问入口
	conf.Runtime.Gin.Handle("GET", "/", index)

	// 静态文件路由
	conf.Runtime.Gin.Static("/css", "./static/css")
	conf.Runtime.Gin.Static("/js", "./static/js")
	conf.Runtime.Gin.Static("/img", "./static/img")
	conf.Runtime.Gin.Static("/fonts", "./static/fonts")
}

func index(c *gin.Context)  {
	c.HTML(200, "index.html", net.Interface{})
}
