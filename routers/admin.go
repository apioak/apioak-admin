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
			service.GET("/name/list", admin.ServiceNameList)
			service.GET("/info/:res_id", admin.ServiceInfo)
			service.PUT("/update/:res_id", admin.ServiceUpdate)
			service.DELETE("/delete/:res_id", admin.ServiceDelete)
			service.PUT("/update/name/:res_id", admin.ServiceUpdateName)
			service.PUT("/switch/enable/:res_id", admin.ServiceSwitchEnable)
			service.PUT("/switch/release/:res_id", admin.ServiceSwitchRelease)
		}

		servicePlugin := adminRouter.Group("service/plugin/config")
		{
			servicePlugin.POST("/add", admin.ServicePluginConfigAdd)
			servicePlugin.GET("/list/:service_res_id", admin.ServicePluginConfigList)
			servicePlugin.GET("/info/:res_id", admin.ServicePluginConfigInfo)
			servicePlugin.PUT("/update/:res_id", admin.ServicePluginConfigUpdate)
			servicePlugin.DELETE("/delete/:res_id", admin.ServicePluginConfigDelete)
			servicePlugin.PUT("/switch/enable/:res_id", admin.ServicePluginConfigSwitchEnable)
		}

		// router
		router := adminRouter.Group("router")
		{
			// router
			router.POST("/add", admin.RouterAdd)
			router.GET("/list", admin.RouterList)
			router.GET("/info/:service_res_id/:router_res_id", admin.RouterInfo)
			router.PUT("/update/:service_res_id/:router_res_id", admin.RouterUpdate)
			router.DELETE("/delete/:service_res_id/:router_res_id", admin.RouterDelete)
			router.PUT("/update/name/:service_res_id/:router_res_id", admin.RouterUpdateName)
			router.PUT("/switch/enable/:service_res_id/:router_res_id", admin.RouterSwitchEnable)
			router.PUT("/switch/release/:service_res_id/:router_res_id", admin.RouterSwitchRelease)
			router.POST("/copy/:service_res_id/:router_res_id", admin.RouterCopy)
		}

		// router plugin
		routerPlugin := adminRouter.Group("router/plugin/config")
		{
			// router plugin
			routerPlugin.POST("/add", admin.RouterPluginConfigAdd)
			routerPlugin.GET("/list/:router_res_id", admin.RouterPluginConfigList)
			routerPlugin.GET("/info/:res_id", admin.RouterPluginConfigInfo)
			routerPlugin.PUT("/update/:res_id", admin.RouterPluginConfigUpdate)
			routerPlugin.DELETE("/delete/:res_id", admin.RouterPluginConfigDelete)
			routerPlugin.PUT("/switch/enable/:res_id", admin.RouterPluginConfigSwitchEnable)
		}

		// upstream
		upstream := adminRouter.Group("upstream")
		{
			upstream.POST("/add", admin.UpstreamAdd)
			upstream.GET("/list", admin.UpstreamList)
			upstream.GET("/info/:res_id", admin.UpstreamInfo)
			upstream.GET("/name/list", admin.UpstreamNameList)
			upstream.PUT("/update/:res_id", admin.UpstreamUpdate)
			upstream.DELETE("/delete/:res_id", admin.UpstreamDelete)
			upstream.PUT("/update/name/:res_id", admin.UpstreamUpdateName)
			upstream.PUT("/switch/enable/:res_id", admin.UpstreamSwitchEnable)
			upstream.PUT("/switch/release/:res_id", admin.UpstreamSwitchRelease)
		}

		// plugin
		plugin := adminRouter.Group("plugin")
		{
			plugin.GET("/type-list", admin.PluginTypeList)
			plugin.GET("/add-list", admin.PluginAddList)
			plugin.GET("/info/:plugin_res_id", admin.PluginInfo)
			// plugin.GET("/list", admin.PluginList)
			// plugin.PUT("/update/:id", admin.PluginUpdate)
			// plugin.DELETE("/delete/:id", admin.PluginDelete)
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
