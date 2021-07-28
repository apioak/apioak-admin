package routes

import "github.com/gin-gonic/gin"

func AdminRegister(routeEngine *gin.Engine) {

	// API——后台管理
	admin := routeEngine.Group("admin")
	{
		// API——用户
		user := admin.Group("user")
		{
			user.GET("/login")
		}

		// API——服务
		service := admin.Group("service")
		{
			service.GET("/info")
		}

		// API——路由
		route := admin.Group("route")
		{
			route.GET("/info")
		}

		// API——插件
		plugin := admin.Group("plugin")
		{
			plugin.GET("/info")
		}

		// API——证书
		certificate := admin.Group("certificate")
		{
			certificate.GET("/info")
		}

		// API——结点
		node := admin.Group("node")
		{
			node.GET("info")
		}
	}
}
