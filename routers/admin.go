package routers

import (
	"apioak-admin/app/controllers/admin"
	"apioak-admin/app/middlewares"
	"github.com/gin-gonic/gin"
)

func RouterRegister(routerEngine *gin.Engine) {

	noLoginRouter := routerEngine.Group("admin")
	{
		user := noLoginRouter.Group("user")
		{
			user.POST("/register", admin.UserRegister)
			user.POST("/login", admin.UserLogin)
		}
	}

	adminRouter := routerEngine.Group("admin", middlewares.CheckUserLogin)
	{
		// user
		user := adminRouter.Group("user")
		{
			user.POST("/logout", admin.UserLogout)
		}

		// service
		service := adminRouter.Group("service")
		{
			service.POST("/add", admin.ServiceAdd)
			service.GET("/list", admin.ServiceList)
			service.GET("/info/:id", admin.ServiceInfo)
			service.PUT("/update/:id", admin.ServiceUpdate)
			service.DELETE("/delete/:id", admin.ServiceDelete)
			service.PUT("/update/name/:id", admin.ServiceUpdateName)
			service.PUT("/switch/enable/:id", admin.ServiceSwitchEnable)
			service.PUT("/switch/release/:id", admin.ServiceSwitchRelease)
		}

		servicePlugin := adminRouter.Group("service/plugin")
		{
			servicePlugin.POST("/add", admin.ServicePluginConfigAdd)
			servicePlugin.GET("/list/:service_id", admin.ServicePluginConfigList)
			servicePlugin.GET("/info/:plugin_config_res_id", admin.ServicePluginConfigInfo)
			servicePlugin.PUT("/update/:plugin_config_res_id", admin.ServicePluginConfigUpdate)
			servicePlugin.DELETE("/delete/:plugin_config_res_id", admin.ServicePluginConfigDelete)
			servicePlugin.PUT("/switch/enable/:plugin_config_res_id", admin.ServicePluginConfigSwitchEnable)
		}

		// router
		route := adminRouter.Group("router")
		{
			// route
			route.POST("/add", admin.RouterAdd)
			route.GET("/list", admin.RouterList)
			route.PUT("/update/:service_res_id/:router_res_id", admin.RouterUpdate)
			route.GET("/info/:service_res_id/:router_res_id", admin.RouterInfo)
			route.PUT("/switch/release/:service_res_id/:router_res_id", admin.RouterSwitchRelease)

			// route.PUT("/update/:service_id/:route_id", admin.RouterUpdate)
			// route.DELETE("/delete/:service_id/:route_id", admin.RouteDelete)
			// route.POST("/copy/:service_id/:source_route_id", admin.RouteCopy)
			// route.PUT("/update/name/:service_id/:route_id", admin.RouteUpdateName)
			// route.PUT("/switch/enable/:service_id/:route_id", admin.RouteSwitchEnable)

			// router plugin
			// route.GET("/plugin/list/:service_id/:route_id", admin.RouterPluginList)
			// route.GET("/plugin/info/:route_id/:plugin_id/:route_plugin_id", admin.RouterPluginInfo)
			// route.POST("/plugin/add/:service_id/:route_id/:plugin_id", admin.RouterPluginAdd)
			// route.PUT("/plugin/update/:route_id/:plugin_id/:route_plugin_id", admin.RouterPluginUpdate)
			// route.DELETE("/plugin/delete/:route_id/:plugin_id/:route_plugin_id", admin.RouterPluginDelete)
			// route.PUT("/plugin/switch/enable/:route_id/:plugin_id/:route_plugin_id", admin.RouterPluginSwitchEnable)
			// route.PUT("/plugin/switch/release/:route_id/:plugin_id/:route_plugin_id", admin.RouterPluginSwitchRelease)
		}

		// plugin
		plugin := adminRouter.Group("plugin")
		{
			plugin.GET("/type-list", admin.PluginTypeList)
			// plugin.GET("/list", admin.PluginList)
			// plugin.GET("/info/:id", admin.PluginInfo)
			// plugin.PUT("/update/:id", admin.PluginUpdate)
			// plugin.DELETE("/delete/:id", admin.PluginDelete)
			plugin.GET("/add-list", admin.PluginAddList)
		}

		// certificate
		certificate := adminRouter.Group("certificate")
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
		clusterNode := adminRouter.Group("cluster-node")
		{
			clusterNode.POST("/add", admin.ClusterNodeAdd)
			clusterNode.GET("/list", admin.ClusterNodeList)
			clusterNode.DELETE("/delete/:id", admin.ClusterNodeDelete)
		}
	}
}
