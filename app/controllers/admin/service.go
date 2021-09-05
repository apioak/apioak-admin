package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func ServiceAdd(c *gin.Context) {

	var serviceAddUpdateStruct = validators.ServiceAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &serviceAddUpdateStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckExistDomain(serviceAddUpdateStruct.ServiceDomains, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceAddUpdateStruct.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateStruct.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateStruct.ServiceDomains)
	serviceNodes := validators.GetServiceAddNodes(serviceAddUpdateStruct.ServiceNodes)

	createErr := services.ServiceCreate(&serviceAddUpdateStruct, &serviceDomains, &serviceNodes)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceInfo(c *gin.Context) {
	serviceId := c.Param("id")

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	structServiceInfo := services.StructServiceInfo{}
	serviceDomainNodeInfo, err := structServiceInfo.ServiceInfoById(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, serviceDomainNodeInfo)
}

func ServiceList(c *gin.Context) {
	var serviceListStruct = validators.ServiceList{}
	if msg, err := packages.ParseRequestParams(c, &serviceListStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	structServiceList := services.StructServiceList{}
	serviceList, total, err := structServiceList.ServiceListPage(&serviceListStruct)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = serviceListStruct
	result.Page = serviceListStruct.Page
	result.PageSize = serviceListStruct.PageSize
	result.Total = total
	result.Data = serviceList

	utils.Ok(c, result)
}

func ServiceUpdate(c *gin.Context) {
	serviceId := c.Param("id")

	var serviceAddUpdateStruct = validators.ServiceAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &serviceAddUpdateStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	err := services.CheckExistDomain(serviceAddUpdateStruct.ServiceDomains, []string{serviceId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceAddUpdateStruct.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateStruct.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateStruct.ServiceDomains)
	serviceNodes := validators.GetServiceAddNodes(serviceAddUpdateStruct.ServiceNodes)

	updateErr := services.ServiceUpdate(serviceId, &serviceAddUpdateStruct, &serviceDomains, &serviceNodes)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceDelete(c *gin.Context) {
	serviceId := c.Param("id")

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	deleteErr := serviceModel.ServiceDelete(serviceId)
	if deleteErr != nil {
		utils.Error(c, deleteErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceUpdateName(c *gin.Context) {
	serviceId := c.Param("id")

	var serviceUpdateNameValidator = validators.ServiceUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &serviceUpdateNameValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	updateNameErr := serviceModel.ServiceUpdateName(serviceId, serviceUpdateNameValidator.Name)
	if updateNameErr != nil {
		utils.Error(c, updateNameErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchEnable(c *gin.Context) {
	serviceId := c.Param("id")

	var serviceSwitchEnableValidator = validators.ServiceSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &serviceSwitchEnableValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	if serviceSwitchEnableValidator.IsEnable == serviceInfo.IsEnable {
		utils.Error(c, enums.CodeMessages(enums.SwitchNoChange))
		return
	}

	updateEnableErr := serviceModel.ServiceSwitchEnable(serviceId, serviceSwitchEnableValidator.IsEnable)
	if updateEnableErr != nil {
		utils.Error(c, updateEnableErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchWebsocket(c *gin.Context) {
	serviceId := c.Param("id")

	var serviceSwitchWebsocketValidator = validators.ServiceSwitchWebsocket{}
	if msg, err := packages.ParseRequestParams(c, &serviceSwitchWebsocketValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	if serviceSwitchWebsocketValidator.WebSocket == serviceInfo.WebSocket {
		utils.Error(c, enums.CodeMessages(enums.SwitchNoChange))
		return
	}

	updateEnableErr := serviceModel.ServiceSwitchWebsocket(serviceId, serviceSwitchWebsocketValidator.WebSocket)
	if updateEnableErr != nil {
		utils.Error(c, updateEnableErr.Error())
		return
	}

	utils.Ok(c)
}
