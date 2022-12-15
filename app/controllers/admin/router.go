package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func RouterAdd(c *gin.Context) {
	bindParams := validators.ValidatorRouterAddUpdate{
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.GetRouterAttributesDefault(&bindParams)
	validators.CorrectUpstreamDefault(&bindParams.UpstreamAddUpdate)
	validators.CorrectUpstreamAddNodes(&bindParams.UpstreamNodes)

	if bindParams.RouterPath == utils.DefaultRouterPath {
		utils.Error(c, enums.CodeMessages(enums.RouterDefaultPathNoPermission))
		return
	}

	checkServiceExistErr := services.CheckServiceExist(bindParams.ServiceResID)
	if (len(bindParams.ServiceResID) == 0) || (checkServiceExistErr != nil) {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	err := services.CheckExistServiceRouterPath(bindParams.ServiceResID, bindParams.RouterPath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	createErr := services.RouterCreate(&bindParams)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterList(c *gin.Context) {
	var bindParams = validators.ValidatorRouterList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(bindParams.ServiceResID)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	structRouterList := services.RouterListItem{}
	routerList, total, err := structRouterList.RouterListPage(bindParams.ServiceResID, &bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = bindParams
	result.Page = bindParams.Page
	result.PageSize = bindParams.PageSize
	result.Total = total
	result.Data = routerList

	utils.Ok(c, result)
}

func RouterInfo(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	structRouterInfo := services.StructRouterInfo{}
	routeInfo, routerInfoErr := structRouterInfo.RouterInfoByServiceRouterId(serviceResId, routerResId)
	if routerInfoErr != nil {
		utils.Error(c, enums.CodeMessages(enums.RouterNull))
		return
	}

	utils.Ok(c, routeInfo)
}

func RouterUpdate(c *gin.Context) {
	var bindParams = validators.ValidatorRouterAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouterAttributesDefault(&bindParams)
	validators.CorrectUpstreamDefault(&bindParams.UpstreamAddUpdate)
	validators.CorrectUpstreamAddNodes(&bindParams.UpstreamNodes)

	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	// @todo 默认服务下的全量路由检测，后续版本中放开
	// checkEditDefaultPathRouterErr := services.CheckEditDefaultPathRouter(routerResId)
	// if checkEditDefaultPathRouterErr != nil {
	// 	utils.Error(c, checkEditDefaultPathRouterErr.Error())
	// 	return
	// }

	checkServiceRouterPathErr := services.CheckServiceRouterPath(bindParams.RouterPath)
	if checkServiceRouterPathErr != nil {
		utils.Error(c, checkServiceRouterPathErr.Error())
		return
	}

	err := services.CheckExistServiceRouterPath(bindParams.ServiceResID, bindParams.RouterPath, []string{serviceResId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := services.RouterUpdate(routerResId, bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterUpdateName(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	var bindParams = validators.RouterUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	routerModel := models.Routers{}
	updateErr := routerModel.RouterUpdateName(routerResId, bindParams.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterSwitchEnable(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routeResId := strings.TrimSpace(c.Param("route_res_id"))

	var bindParams = validators.RouterSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routeResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRouter(routeResId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	checkRouteEnableChangeErr := services.CheckRouterEnableChange(routeResId, bindParams.Enable)
	if checkRouteEnableChangeErr != nil {
		utils.Error(c, checkRouteEnableChangeErr.Error())
		return
	}

	routerModel := models.Routers{}
	updateErr := routerModel.RouterSwitchEnable(routeResId, bindParams.Enable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}


func RouterSwitchRelease(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	// @todo 这里增加检测服务是否已发布，只要不是未发布都可以通过，请优先发布服务（未发布的服务不允许发布路由）

	checkExistRouterErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouterErr != nil {
		utils.Error(c, checkExistRouterErr.Error())
		return
	}

	checkRouterReleaseErr := services.CheckRouterRelease(routerResId)
	if checkRouterReleaseErr != nil {
		utils.Error(c, checkRouterReleaseErr.Error())
		return
	}

	serviceRouterReleaseErr := services.RouterUpstreamRelease([]string{routerResId}, utils.ReleaseTypePush)
	if serviceRouterReleaseErr != nil {
		utils.Error(c, serviceRouterReleaseErr.Error())
		return
	}

	utils.Ok(c)
}

// ------------------------------------------------------------------------------------------



// func RouteDelete(c *gin.Context) {
// 	serviceId := strings.TrimSpace(c.Param("service_id"))
// 	routeId := strings.TrimSpace(c.Param("route_id"))
//
// 	checkServiceExistErr := services.CheckServiceExist(serviceId)
// 	if checkServiceExistErr != nil {
// 		utils.Error(c, checkServiceExistErr.Error())
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRoute(routeId)
// 	if checkEditDefaultPathRouteErr != nil {
// 		utils.Error(c, checkEditDefaultPathRouteErr.Error())
// 		return
// 	}
//
// 	checkRouteDeleteErr := services.CheckRouteDelete(routeId)
// 	if checkRouteDeleteErr != nil {
// 		utils.Error(c, checkRouteDeleteErr.Error())
// 		return
// 	}
//
// 	deleteErr := services.RouteDelete(routeId)
// 	if deleteErr != nil {
// 		utils.Error(c, enums.CodeMessages(enums.Error))
// 		return
// 	}
//
// 	utils.Ok(c)
// }

// func RouteCopy(c *gin.Context) {
// 	var bindParams = validators.ValidatorRouterAddUpdate{
// 		Release: utils.ReleaseN,
// 	}
// 	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
// 		utils.Error(c, msg)
// 		return
// 	}
// 	validators.GetRouterAttributesDefault(&bindParams)
//
// 	if bindParams.RoutePath == utils.DefaultRoutePath {
// 		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
// 		return
// 	}
//
// 	serviceId := strings.TrimSpace(c.Param("service_id"))
// 	sourceRouteId := strings.TrimSpace(c.Param("source_route_id"))
//
// 	checkServiceExistErr := services.CheckServiceExist(serviceId)
// 	if checkServiceExistErr != nil {
// 		utils.Error(c, checkServiceExistErr.Error())
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(sourceRouteId, serviceId)
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	err := services.CheckExistServiceRoutePath(bindParams.ServiceResID, bindParams.RoutePath, []string{})
// 	if err != nil {
// 		utils.Error(c, err.Error())
// 		return
// 	}
//
// 	createErr := services.RouteCopy(&bindParams, sourceRouteId)
// 	if createErr != nil {
// 		utils.Error(c, createErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c)
// }



// func RoutePluginList(c *gin.Context) {
// 	serviceId := strings.TrimSpace(c.Param("service_id"))
// 	routeId := strings.TrimSpace(c.Param("route_id"))
//
// 	checkServiceExistErr := services.CheckServiceExist(serviceId)
// 	if checkServiceExistErr != nil {
// 		utils.Error(c, checkServiceExistErr.Error())
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	routePluginInfo := services.RoutePluginInfo{}
// 	routePluginInfoList := routePluginInfo.RoutePluginList(routeId)
//
// 	utils.Ok(c, routePluginInfoList)
// }
//
// func RoutePluginAdd(c *gin.Context) {
// 	serviceId := strings.TrimSpace(c.Param("service_id"))
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
//
// 	var bindParams = validators.RoutePluginAddUpdate{}
// 	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
// 		utils.Error(c, msg)
// 		return
// 	}
// 	bindParams.RouteID = routeId
// 	bindParams.PluginID = pluginId
//
// 	checkServiceExistErr := services.CheckServiceExist(serviceId)
// 	if checkServiceExistErr != nil {
// 		utils.Error(c, checkServiceExistErr.Error())
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExistErr := services.CheckRoutePluginExistByRoutePluginId(routeId, pluginId)
// 	if checkRoutePluginExistErr != nil {
// 		utils.Error(c, checkRoutePluginExistErr.Error())
// 		return
// 	}
//
// 	checkPluginConfigErr := services.CheckPluginConfig(pluginId, &bindParams)
// 	if checkPluginConfigErr != nil {
// 		utils.Error(c, checkPluginConfigErr.Error())
// 		return
// 	}
//
// 	addErr := services.RoutePluginCreate(&bindParams)
// 	if addErr != nil {
// 		utils.Error(c, addErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c)
// }
//
// func RoutePluginUpdate(c *gin.Context) {
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
// 	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))
//
// 	var bindParams = validators.RoutePluginAddUpdate{}
// 	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
// 		utils.Error(c, msg)
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, "")
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
// 	if checkRoutePluginExist != nil {
// 		utils.Error(c, checkRoutePluginExist.Error())
// 		return
// 	}
//
// 	checkPluginConfigErr := services.CheckPluginConfig(pluginId, &bindParams)
// 	if checkPluginConfigErr != nil {
// 		utils.Error(c, checkPluginConfigErr.Error())
// 		return
// 	}
//
// 	updateErr := services.RoutePluginUpdate(routePluginId, &bindParams)
// 	if updateErr != nil {
// 		utils.Error(c, updateErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c)
// }
//
// func RoutePluginDelete(c *gin.Context) {
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
// 	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, "")
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
// 	if checkRoutePluginExist != nil {
// 		utils.Error(c, checkRoutePluginExist.Error())
// 		return
// 	}
//
// 	checkRoutePluginDeleteErr := services.CheckRoutePluginDelete(routePluginId)
// 	if checkRoutePluginDeleteErr != nil {
// 		utils.Error(c, checkRoutePluginDeleteErr.Error())
// 		return
// 	}
//
// 	deleteErr := services.RoutePluginDelete(routePluginId)
// 	if deleteErr != nil {
// 		utils.Error(c, enums.CodeMessages(enums.Error))
// 		return
// 	}
//
// 	utils.Ok(c)
// }
//
// func RoutePluginSwitchEnable(c *gin.Context) {
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
// 	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))
//
// 	var bindParams = validators.RoutePluginSwitchEnable{}
// 	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
// 		utils.Error(c, msg)
// 		return
// 	}
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, "")
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
// 	if checkRoutePluginExist != nil {
// 		utils.Error(c, checkRoutePluginExist.Error())
// 		return
// 	}
//
// 	checkRoutePluginEnableChange := services.CheckRoutePluginEnableChange(routePluginId, bindParams.IsEnable)
// 	if checkRoutePluginEnableChange != nil {
// 		utils.Error(c, checkRoutePluginEnableChange.Error())
// 		return
// 	}
//
// 	routePluginModel := models.RoutePlugins{}
// 	updateErr := routePluginModel.RoutePluginSwitchEnable(routePluginId, bindParams.IsEnable)
// 	if updateErr != nil {
// 		utils.Error(c, updateErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c)
// }
//
// func RoutePluginSwitchRelease(c *gin.Context) {
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
// 	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, "")
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
// 	if checkRoutePluginExist != nil {
// 		utils.Error(c, checkRoutePluginExist.Error())
// 		return
// 	}
//
// 	checkRoutePluginReleaseErr := services.CheckRoutePluginRelease(routePluginId)
// 	if checkRoutePluginReleaseErr != nil {
// 		utils.Error(c, checkRoutePluginReleaseErr.Error())
// 		return
// 	}
//
// 	routePluginReleaseErr := services.RoutePluginRelease(routePluginId)
// 	if routePluginReleaseErr != nil {
// 		utils.Error(c, routePluginReleaseErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c)
// }
//
// func RoutePluginInfo(c *gin.Context) {
// 	routeId := strings.TrimSpace(c.Param("route_id"))
// 	pluginId := strings.TrimSpace(c.Param("plugin_id"))
// 	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))
//
// 	checkExistRouteErr := services.CheckRouteExist(routeId, "")
// 	if checkExistRouteErr != nil {
// 		utils.Error(c, checkExistRouteErr.Error())
// 		return
// 	}
//
// 	checkPluginExistErr := services.CheckPluginExist(pluginId)
// 	if checkPluginExistErr != nil {
// 		utils.Error(c, checkPluginExistErr.Error())
// 		return
// 	}
//
// 	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
// 	if checkRoutePluginExist != nil {
// 		utils.Error(c, checkRoutePluginExist.Error())
// 		return
// 	}
//
// 	routePluginConfigInfo, routePluginConfigInfoErr := services.RoutePluginConfigInfo(routePluginId)
// 	if routePluginConfigInfoErr != nil {
// 		utils.Error(c, routePluginConfigInfoErr.Error())
// 		return
// 	}
//
// 	utils.Ok(c, routePluginConfigInfo)
// }
