package routes

import (
	"apioak-admin/app/controllers/admin"
	"apioak-admin/app/middlewares"
	"github.com/gin-gonic/gin"
)

func RouteRegister(routeEngine *gin.Engine) {

	noLoginRoute := routeEngine.Group("admin")
	{
		user := noLoginRoute.Group("user")
		{
			user.POST("/register", admin.UserRegister)
			user.POST("/login", admin.UserLogin)
		}
	}

	adminRoute := routeEngine.Group("admin", middlewares.CheckUserLogin)
	{
		// user
		user := adminRoute.Group("user")
		{
			user.POST("/logout", admin.UserLogout)
		}

		// service
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
			service.PUT("/switch/release/:id", admin.ServiceSwitchRelease)
			service.PUT("/switch/websocket/:id", admin.ServiceSwitchWebsocket)
			service.PUT("/switch/health-check/:id", admin.ServiceSwitchHealthCheck)
		}

		// route
		route := adminRoute.Group("route")
		{
			// route
			route.POST("/add/:service_id", admin.RouteAdd)
			route.GET("/list/:service_id", admin.RouteList)
			route.GET("/info/:service_id/:route_id", admin.RouteInfo)
			route.PUT("/update/:service_id/:route_id", admin.RouteUpdate)
			route.DELETE("/delete/:service_id/:route_id", admin.RouteDelete)
			route.POST("/copy/:service_id/:source_route_id", admin.RouteCopy)
			route.PUT("/update/name/:service_id/:route_id", admin.RouteUpdateName)
			route.PUT("/switch/enable/:service_id/:route_id", admin.RouteSwitchEnable)
			route.PUT("/switch/release/:service_id/:route_id", admin.RouteSwitchRelease)

			// route plugin
			route.GET("/add-plugin/list/:service_id/:route_id", admin.RoutePluginFilterList)
			route.GET("/plugin/list/:service_id/:route_id", admin.RoutePluginList)
			route.GET("/plugin/info/:route_id/:plugin_id/:route_plugin_id", admin.RoutePluginInfo)
			route.POST("/plugin/add/:service_id/:route_id/:plugin_id", admin.RoutePluginAdd)
			route.PUT("/plugin/update/:route_id/:plugin_id/:route_plugin_id", admin.RoutePluginUpdate)
			route.DELETE("/plugin/delete/:route_id/:plugin_id/:route_plugin_id", admin.RoutePluginDelete)
			route.PUT("/plugin/switch/enable/:route_id/:plugin_id/:route_plugin_id", admin.RoutePluginSwitchEnable)
			route.PUT("/plugin/switch/release/:route_id/:plugin_id/:route_plugin_id", admin.RoutePluginSwitchRelease)
		}

		// plugin
		plugin := adminRoute.Group("plugin")
		{
			plugin.GET("/type-list", admin.PluginTypeList)
			plugin.POST("/add", admin.PluginAdd)
			plugin.GET("/list", admin.PluginList)
			plugin.GET("/info/:id", admin.PluginInfo)
			plugin.PUT("/update/:id", admin.PluginUpdate)
			plugin.DELETE("/delete/:id", admin.PluginDelete)
		}

		// certificate
		certificate := adminRoute.Group("certificate")
		{
			certificate.GET("/list", admin.CertificateList)
			certificate.POST("/add", admin.CertificateAdd)
			certificate.GET("/info/:id", admin.CertificateInfo)
			certificate.PUT("/update/:id", admin.CertificateUpdate)
			certificate.DELETE("/delete/:id", admin.CertificateDelete)
			certificate.PUT("/switch/enable/:id", admin.CertificateSwitchEnable)
			certificate.PUT("/switch/release/:id", admin.CertificateSwitchRelease)
		}

		// cluster node
		clusterNode := adminRoute.Group("cluster-node")
		{
			clusterNode.GET("list", admin.ClusterNodeList)
			clusterNode.DELETE("/delete/:id", admin.ClusterNodeDelete)
			clusterNode.PUT("/switch/enable/:id", admin.ClusterNodeSwitchEnable)
		}
	}
}