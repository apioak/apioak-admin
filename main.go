package main

import (
	"apioak-admin/cores"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
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

/**
 * 初始化静态文件服务
 */
func initStatic(conf *cores.ConfigGlobal) {

	// 引入静态文件模板
	conf.Runtime.Gin.SetHTMLTemplate(template.Must(template.New("").ParseFS(htmlFS, "html/*.html")))

	// 静态文件系统目录
	assetsFs, _ := fs.Sub(htmlFS, "html/assets")

	// 路由匹配并加载静态文件系统
	conf.Runtime.Gin.StaticFS("/assets", http.FS(assetsFs))

	// 访问入口
	conf.Runtime.Gin.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
}
