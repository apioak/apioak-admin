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
			service.GET("/common/load-balance/list", admin.ServiceLoadBalanceList)

			service.POST("/add", admin.ServiceAdd)
			service.GET("/list", admin.ServiceList)
			service.GET("/info/:id", admin.ServiceInfo)
			service.PUT("/update/:id", admin.ServiceUpdate)
			service.DELETE("/delete/:id", admin.ServiceDelete)
			service.PUT("/update/name/:id", admin.ServiceUpdateName)
			service.PUT("/switch/enable/:id", admin.ServiceSwitchEnable)
			service.PUT("/switch/websocket/:id", admin.ServiceSwitchWebsocket)
			service.PUT("/switch/health-check/:id", admin.ServiceSwitchHealthCheck)
		}

		// 路由
		route := adminRoute.Group("route")
		{
			route.POST("/add", admin.RouteAdd)
			route.GET("/list/:service_id", admin.RouteList)
			//route.GET("/info/:id", admin.RouteInfo)
			//route.PUT("/update/:id", admin.RouteUpdate)
			//route.DELETE("/delete/:id", admin.RouteDelete)
			//route.PUT("/update/name/:id/", admin.RouteUpdateName)
			//route.PUT("/switch/enable/:id", admin.RouteSwitchEnable)
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
