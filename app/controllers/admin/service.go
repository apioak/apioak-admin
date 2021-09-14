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

func ServiceLoadBalanceList(c *gin.Context) {
	loadBalance := utils.LoadBalance{}
	loadBalanceList := loadBalance.LoadBalanceList()
	utils.Ok(c, loadBalanceList)
}

func ServiceAdd(c *gin.Context) {

	var serviceAddUpdateValidator = validators.ServiceAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &serviceAddUpdateValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckExistDomain(serviceAddUpdateValidator.ServiceDomains, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	validators.GetServiceAttributesDefault(&serviceAddUpdateValidator)
	serviceAddUpdateValidator.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateValidator.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateValidator.ServiceDomains)
	serviceNodes := validators.GetServiceAddNodes(serviceAddUpdateValidator.ServiceNodes)

	createErr := services.ServiceCreate(&serviceAddUpdateValidator, &serviceDomains, &serviceNodes)
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

	var serviceAddUpdateValidator = validators.ServiceAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &serviceAddUpdateValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	err := services.CheckExistDomain(serviceAddUpdateValidator.ServiceDomains, []string{serviceId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceAddUpdateValidator.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateValidator.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateValidator.ServiceDomains)
	serviceNodes := validators.GetServiceAddNodes(serviceAddUpdateValidator.ServiceNodes)

	updateErr := services.ServiceUpdate(serviceId, &serviceAddUpdateValidator, &serviceDomains, &serviceNodes)
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

	if serviceInfo.IsEnable == utils.EnableOn {
		utils.Error(c, enums.CodeMessages(enums.SwitchONProhibitsOp))
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

	updateErr := serviceModel.ServiceUpdateName(serviceId, serviceUpdateNameValidator.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
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

	updateErr := serviceModel.ServiceSwitchEnable(serviceId, serviceSwitchEnableValidator.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
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

	updateErr := serviceModel.ServiceSwitchWebsocket(serviceId, serviceSwitchWebsocketValidator.WebSocket)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchHealthCheck(c *gin.Context) {
	serviceId := c.Param("id")

	var serviceSwitchHealthCheckValidator = validators.ServiceSwitchHealthCheck{}
	if msg, err := packages.ParseRequestParams(c, &serviceSwitchHealthCheckValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	if serviceSwitchHealthCheckValidator.HealthCheck == serviceInfo.HealthCheck {
		utils.Error(c, enums.CodeMessages(enums.SwitchNoChange))
		return
	}

	updateErr := serviceModel.ServiceSwitchHealthCheck(serviceId, serviceSwitchHealthCheckValidator.HealthCheck)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}
