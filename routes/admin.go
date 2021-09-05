package routes

import (
	"apioak-admin/app/controllers/admin"
	"apioak-admin/app/middlewares"
	"github.com/gin-gonic/gin"
)

func RouteRegister(routeEngine *gin.Engine) {

	adminRoute := routeEngine.Group("admin")
	{
		// 用户
		user := adminRoute.Group("user").Use(middlewares.CheckUserLogin)
		{
			user.POST("/register", admin.UserRegister)
		}

		// 服务
		service := adminRoute.Group("service")
		{
			service.GET("/list", admin.ServiceList)
			//service.GET("/info/:id", admin.ServiceInfo)
			service.POST("/add", admin.ServiceAdd)
			service.PUT("/update/:id", admin.ServiceUpdate)
			service.DELETE("/delete/:id", admin.ServiceDelete)
			//service.PUT("/name/update/:id/", admin.ServiceDelete)
			//service.PUT("/switch/health-check/:id/:enable", admin.ServiceDelete)
			//service.PUT("/switch/websocket/:id/:enable", admin.ServiceDelete)
			//service.PUT("/switch/enable/:id/:enable", admin.ServiceDelete)
		}

		// 路由
		route := adminRoute.Group("route")
		{
			route.GET("/info")
		}

		// 插件
		plugin := adminRoute.Group("plugin")
		{
			plugin.GET("/info")
		}

		// 证书
		certificate := adminRoute.Group("certificate")
		{
			certificate.GET("/info")
		}

		// 节点
		node := adminRoute.Group("node")
		{
			node.GET("info")
		}
	}
}
