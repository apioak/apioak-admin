package routes

import (
	"apioak-admin/app/controllers/admin"
	"github.com/gin-gonic/gin"
)

func RouteRegister(routeEngine *gin.Engine) {

	// API——后台管理
	adminRoute := routeEngine.Group("admin")
	{
		// 用户
		user := adminRoute.Group("user")
		{
			user.POST("/register", admin.Register)
		}

		// 服务
		service := adminRoute.Group("service")
		{
			service.GET("/info")
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
