package routes

import (
	"apioak-admin/app/controllers/admin"
	"apioak-admin/app/middlewares"
	"github.com/gin-gonic/gin"
)

func RouteRegister(routeEngine *gin.Engine) {

	adminRoute := routeEngine.Group("admin", middlewares.CheckUserLogin)
	{
		// 用户
		user := adminRoute.Group("user")
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
			route.GET("/info/:service_id/:id", admin.RouteInfo)
			route.PUT("/update/:service_id/:id", admin.RouteUpdate)
			route.DELETE("/delete/:service_id/:id", admin.RouteDelete)
			route.PUT("/update/name/:service_id/:id", admin.RouteUpdateName)
			route.PUT("/switch/enable/:service_id/:id", admin.RouteSwitchEnable)

			route.GET("/add-plugin/list/:service_id/:id", admin.RoutePluginFilterList)
			//route.GET("/plugin/list/:id", admin.RoutePluginList)
			//route.POST("/plugin/add/:id/:plugin_id", admin.RoutePluginAdd)
			//route.GET("/plugin/info/:id/:plugin_id", admin.RoutePluginInfo)
			//route.PUT("/plugin/update/:id/:plugin_id", admin.RoutePluginUpdate)
			//route.DELETE("/plugin/delete/:id/:plugin_id", admin.RoutePluginDelete)
			//route.DELETE("/plugin/switch/enable/:id/:plugin_id", admin.RoutePluginSwitchEnable)
		}

		// 插件
		plugin := adminRoute.Group("plugin")
		{
			plugin.POST("/add", admin.PluginAdd)
			plugin.GET("/list", admin.PluginList)
			plugin.PUT("/update/:id", admin.PluginUpdate)
			plugin.DELETE("/delete/:id", admin.PluginDelete)
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
