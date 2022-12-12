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

func RouteAdd(c *gin.Context) {
	bindParams := validators.ValidatorRouteAddUpdate{
		Release: utils.ReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.GetRouteAttributesDefault(&bindParams)
	validators.CorrectUpstreamTimeOut(&bindParams.UpstreamTimeout)
	validators.CorrectUpstreamAddNodes(&bindParams.UpstreamNodes)

	if bindParams.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	// @todo 这里检测方法内需要改动，牵扯到服务了
	checkServiceExistErr := services.CheckServiceExist(bindParams.ServiceResID)
	if (len(bindParams.ServiceResID) == 0) || (checkServiceExistErr != nil) {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	err := services.CheckExistServiceRoutePath(bindParams.ServiceResID, bindParams.RoutePath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	createErr := services.RouteCreate(&bindParams)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteList(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	var bindParams = validators.ValidatorRouteList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	structRouteList := services.StructRouteList{}
	routeList, total, err := structRouteList.RouteListPage(serviceId, &bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = bindParams
	result.Page = bindParams.Page
	result.PageSize = bindParams.PageSize
	result.Total = total
	result.Data = routeList

	utils.Ok(c, result)
}

func RouteInfo(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	structRouteInfo := services.StructRouteInfo{}
	routeInfo, routeInfoErr := structRouteInfo.RouteInfoByServiceRouteId(serviceId, routeId)
	if routeInfoErr != nil {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	utils.Ok(c, routeInfo)
}

func RouteUpdate(c *gin.Context) {
	var bindParams = validators.ValidatorRouteAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouteAttributesDefault(&bindParams)

	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRoute(routeId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	checkServiceRoutePathErr := services.CheckServiceRoutePath(bindParams.RoutePath)
	if checkServiceRoutePathErr != nil {
		utils.Error(c, checkServiceRoutePathErr.Error())
		return
	}

	err := services.CheckExistServiceRoutePath(bindParams.ServiceResID, bindParams.RoutePath, []string{routeId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := services.RouteUpdate(routeId, &bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteDelete(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRoute(routeId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	checkRouteDeleteErr := services.CheckRouteDelete(routeId)
	if checkRouteDeleteErr != nil {
		utils.Error(c, checkRouteDeleteErr.Error())
		return
	}

	deleteErr := services.RouteDelete(routeId)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func RouteCopy(c *gin.Context) {
	var bindParams = validators.ValidatorRouteAddUpdate{
		Release: utils.ReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouteAttributesDefault(&bindParams)

	if bindParams.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	serviceId := strings.TrimSpace(c.Param("service_id"))
	sourceRouteId := strings.TrimSpace(c.Param("source_route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(sourceRouteId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	err := services.CheckExistServiceRoutePath(bindParams.ServiceResID, bindParams.RoutePath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	createErr := services.RouteCopy(&bindParams, sourceRouteId)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteUpdateName(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	var bindParams = validators.RouteUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	routeModel := models.Routes{}
	updateErr := routeModel.RouteUpdateName(routeId, bindParams.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteSwitchEnable(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	var bindParams = validators.RouteSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRoute(routeId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	checkRouteEnableChangeErr := services.CheckRouteEnableChange(routeId, bindParams.IsEnable)
	if checkRouteEnableChangeErr != nil {
		utils.Error(c, checkRouteEnableChangeErr.Error())
		return
	}

	routeModel := models.Routes{}
	updateErr := routeModel.RouteSwitchEnable(routeId, bindParams.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteSwitchRelease(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routeResId := strings.TrimSpace(c.Param("route_res_id"))

	// checkServiceExistErr := services.CheckServiceExist(serviceResId)
	// if checkServiceExistErr != nil {
	// 	utils.Error(c, checkServiceExistErr.Error())
	// 	return
	// }

	// @todo 检测服务是否已发布

	checkExistRouteErr := services.CheckRouteExist(routeResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkRouteReleaseErr := services.CheckRouteRelease(routeResId)
	if checkRouteReleaseErr != nil {
		utils.Error(c, checkRouteReleaseErr.Error())
		return
	}

	serviceRouteReleaseErr := services.ServiceRouteRelease([]string{routeResId}, utils.ReleaseTypePush)
	if serviceRouteReleaseErr != nil {
		utils.Error(c, serviceRouteReleaseErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginFilterList(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	routeAddPluginInfo := services.RouteAddPluginInfo{}
	routeAddPluginList, routeAddPluginListErr := routeAddPluginInfo.RouteAddPluginList(routeId)
	if routeAddPluginListErr != nil {
		utils.Error(c, routeAddPluginListErr.Error())
		return
	}

	utils.Ok(c, routeAddPluginList)
}

func RoutePluginList(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	routePluginInfo := services.RoutePluginInfo{}
	routePluginInfoList := routePluginInfo.RoutePluginList(routeId)

	utils.Ok(c, routePluginInfoList)
}

func RoutePluginAdd(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))

	var bindParams = validators.RoutePluginAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}
	bindParams.RouteID = routeId
	bindParams.PluginID = pluginId

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, serviceId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExistErr := services.CheckRoutePluginExistByRoutePluginId(routeId, pluginId)
	if checkRoutePluginExistErr != nil {
		utils.Error(c, checkRoutePluginExistErr.Error())
		return
	}

	checkPluginConfigErr := services.CheckPluginConfig(pluginId, &bindParams)
	if checkPluginConfigErr != nil {
		utils.Error(c, checkPluginConfigErr.Error())
		return
	}

	addErr := services.RoutePluginCreate(&bindParams)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginUpdate(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))

	var bindParams = validators.RoutePluginAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, "")
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
	if checkRoutePluginExist != nil {
		utils.Error(c, checkRoutePluginExist.Error())
		return
	}

	checkPluginConfigErr := services.CheckPluginConfig(pluginId, &bindParams)
	if checkPluginConfigErr != nil {
		utils.Error(c, checkPluginConfigErr.Error())
		return
	}

	updateErr := services.RoutePluginUpdate(routePluginId, &bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginDelete(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))

	checkExistRouteErr := services.CheckRouteExist(routeId, "")
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
	if checkRoutePluginExist != nil {
		utils.Error(c, checkRoutePluginExist.Error())
		return
	}

	checkRoutePluginDeleteErr := services.CheckRoutePluginDelete(routePluginId)
	if checkRoutePluginDeleteErr != nil {
		utils.Error(c, checkRoutePluginDeleteErr.Error())
		return
	}

	deleteErr := services.RoutePluginDelete(routePluginId)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func RoutePluginSwitchEnable(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))

	var bindParams = validators.RoutePluginSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkExistRouteErr := services.CheckRouteExist(routeId, "")
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
	if checkRoutePluginExist != nil {
		utils.Error(c, checkRoutePluginExist.Error())
		return
	}

	checkRoutePluginEnableChange := services.CheckRoutePluginEnableChange(routePluginId, bindParams.IsEnable)
	if checkRoutePluginEnableChange != nil {
		utils.Error(c, checkRoutePluginEnableChange.Error())
		return
	}

	routePluginModel := models.RoutePlugins{}
	updateErr := routePluginModel.RoutePluginSwitchEnable(routePluginId, bindParams.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginSwitchRelease(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))

	checkExistRouteErr := services.CheckRouteExist(routeId, "")
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
	if checkRoutePluginExist != nil {
		utils.Error(c, checkRoutePluginExist.Error())
		return
	}

	checkRoutePluginReleaseErr := services.CheckRoutePluginRelease(routePluginId)
	if checkRoutePluginReleaseErr != nil {
		utils.Error(c, checkRoutePluginReleaseErr.Error())
		return
	}

	routePluginReleaseErr := services.RoutePluginRelease(routePluginId)
	if routePluginReleaseErr != nil {
		utils.Error(c, routePluginReleaseErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginInfo(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("route_plugin_id"))

	checkExistRouteErr := services.CheckRouteExist(routeId, "")
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkPluginExistErr := services.CheckPluginExist(pluginId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	checkRoutePluginExist := services.CheckRoutePluginExist(routePluginId, routeId, pluginId)
	if checkRoutePluginExist != nil {
		utils.Error(c, checkRoutePluginExist.Error())
		return
	}

	routePluginConfigInfo, routePluginConfigInfoErr := services.RoutePluginConfigInfo(routePluginId)
	if routePluginConfigInfoErr != nil {
		utils.Error(c, routePluginConfigInfoErr.Error())
		return
	}

	utils.Ok(c, routePluginConfigInfo)
}
