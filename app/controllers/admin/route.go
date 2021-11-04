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
	serviceId := strings.TrimSpace(c.Param("service_id"))

	validatorRouteAddUpdate := validators.ValidatorRouteAddUpdate{}
	validatorRouteAddUpdate.ServiceID = serviceId
	if msg, err := packages.ParseRequestParams(c, &validatorRouteAddUpdate); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouteAttributesDefault(&validatorRouteAddUpdate)

	if validatorRouteAddUpdate.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	err := services.CheckExistServiceRoutePath(validatorRouteAddUpdate.ServiceID, validatorRouteAddUpdate.RoutePath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	createErr := services.RouteCreate(&validatorRouteAddUpdate)
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

	var validatorRouteList = validators.ValidatorRouteList{}
	if msg, err := packages.ParseRequestParams(c, &validatorRouteList); err != nil {
		utils.Error(c, msg)
		return
	}

	structRouteList := services.StructRouteList{}
	routeList, total, err := structRouteList.RouteListPage(serviceId, &validatorRouteList)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = validatorRouteList
	result.Page = validatorRouteList.Page
	result.PageSize = validatorRouteList.PageSize
	result.Total = total
	result.Data = routeList

	utils.Ok(c, result)
}

func RouteInfo(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

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
	var validatorRouteAddUpdate = validators.ValidatorRouteAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &validatorRouteAddUpdate); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouteAttributesDefault(&validatorRouteAddUpdate)

	if validatorRouteAddUpdate.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

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

	err := services.CheckExistServiceRoutePath(validatorRouteAddUpdate.ServiceID, validatorRouteAddUpdate.RoutePath, []string{routeId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := services.RouteUpdate(routeId, &validatorRouteAddUpdate)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteDelete(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

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

	checkRouteEnableOnProhibitedOpErr := services.CheckRouteEnableOnProhibitedOp(routeId)
	if checkRouteEnableOnProhibitedOpErr != nil {
		utils.Error(c, checkRouteEnableOnProhibitedOpErr.Error())
		return
	}

	deleteErr := services.RouteDelete(routeId)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func RouteUpdateName(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

	var routeUpdateNameValidator = validators.RouteUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &routeUpdateNameValidator); err != nil {
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
	updateErr := routeModel.RouteUpdateName(routeId, routeUpdateNameValidator.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteSwitchEnable(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

	var routeSwitchEnableValidator = validators.RouteSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &routeSwitchEnableValidator); err != nil {
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

	checkRouteEnableChangeErr := services.CheckRouteEnableChange(routeId, routeSwitchEnableValidator.IsEnable)
	if checkRouteEnableChangeErr != nil {
		utils.Error(c, checkRouteEnableChangeErr.Error())
		return
	}

	updateErr := services.ServiceRouteSwitchEnable(routeId, routeSwitchEnableValidator.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginFilterList(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))
	routeId := strings.TrimSpace(c.Param("id"))

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
	routeId := strings.TrimSpace(c.Param("id"))

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

	var routePluginAddValidator = validators.RoutePluginAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &routePluginAddValidator); err != nil {
		utils.Error(c, msg)
		return
	}
	routePluginAddValidator.RouteID = routeId
	routePluginAddValidator.PluginID = pluginId

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

	checkPluginConfigErr := services.CheckPluginConfig(pluginId, routePluginAddValidator.Config)
	if checkPluginConfigErr != nil {
		utils.Error(c, checkPluginConfigErr.Error())
		return
	}

	addErr := services.RoutePluginCreate(&routePluginAddValidator)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginUpdate(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("id"))

	var routePluginUpdateValidator = validators.RoutePluginAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &routePluginUpdateValidator); err != nil {
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

	checkPluginConfigErr := services.CheckPluginConfig(pluginId, routePluginUpdateValidator.Config)
	if checkPluginConfigErr != nil {
		utils.Error(c, checkPluginConfigErr.Error())
		return
	}

	updateErr := services.RoutePluginUpdate(routePluginId, &routePluginUpdateValidator)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginDelete(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("id"))

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

	checkRoutePluginEnableOn := services.CheckRoutePluginEnableOn(routePluginId)
	if checkRoutePluginEnableOn != nil {
		utils.Error(c, checkRoutePluginEnableOn.Error())
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
	routePluginId := strings.TrimSpace(c.Param("id"))

	var routePluginSwitchEnableValidator = validators.RoutePluginSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &routePluginSwitchEnableValidator); err != nil {
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

	checkRoutePluginEnableChange := services.CheckRoutePluginEnableChange(routePluginId, routePluginSwitchEnableValidator.IsEnable)
	if checkRoutePluginEnableChange != nil {
		utils.Error(c, checkRoutePluginEnableChange.Error())
		return
	}

	updateErr := services.RoutePluginSwitchEnable(routePluginId, routePluginSwitchEnableValidator.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RoutePluginInfo(c *gin.Context) {
	routeId := strings.TrimSpace(c.Param("route_id"))
	pluginId := strings.TrimSpace(c.Param("plugin_id"))
	routePluginId := strings.TrimSpace(c.Param("id"))

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
