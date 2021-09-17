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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(validatorRouteAddUpdate.ServiceID)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosById(routeId)
	if routeModelInfoErr != nil {
		utils.Error(c, routeModelInfoErr.Error())
		return
	}

	if routeModelInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	if routeModelInfo.ServiceID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.RouteServiceNoMatch))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosById(routeId)
	if routeModelInfoErr != nil {
		utils.Error(c, routeModelInfoErr.Error())
		return
	}

	if routeModelInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	if routeModelInfo.ServiceID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.RouteServiceNoMatch))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosById(routeId)
	if routeModelInfoErr != nil {
		utils.Error(c, routeModelInfoErr.Error())
		return
	}

	if routeModelInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	if routeModelInfo.ServiceID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.RouteServiceNoMatch))
		return
	}

	if routeModelInfo.IsEnable == utils.EnableOn {
		utils.Error(c, enums.CodeMessages(enums.SwitchONProhibitsOp))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosById(routeId)
	if routeModelInfoErr != nil {
		utils.Error(c, routeModelInfoErr.Error())
		return
	}

	if routeModelInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	if routeModelInfo.ServiceID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.RouteServiceNoMatch))
		return
	}

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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosById(routeId)
	if routeModelInfoErr != nil {
		utils.Error(c, routeModelInfoErr.Error())
		return
	}

	if routeModelInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
		return
	}

	if routeModelInfo.ServiceID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.RouteServiceNoMatch))
		return
	}

	if routeSwitchEnableValidator.IsEnable == routeModelInfo.IsEnable {
		utils.Error(c, enums.CodeMessages(enums.SwitchNoChange))
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

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	routeModel := &models.Routes{}
	routeInfo, routeInfoErr := routeModel.RouteInfosById(routeId)
	if routeInfoErr != nil {
		utils.Error(c, routeInfoErr.Error())
		return
	}
	if routeInfo.ID != routeId {
		utils.Error(c, enums.CodeMessages(enums.RouteNull))
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
