package admin

import (
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func ServiceLoadBalanceList(c *gin.Context) {
	loadBalanceList := utils.LoadBalanceList()
	utils.Ok(c, loadBalanceList)
}

func ServiceAdd(c *gin.Context) {

	var bindParams = validators.ServiceAddUpdate{
		IsRelease: utils.IsReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.CorrectServiceAttributesDefault(&bindParams)
	validators.CorrectServiceTimeOut(&bindParams.Timeouts)
	validators.CorrectServiceDomains(&bindParams.ServiceDomains)
	validators.CorrectServiceAddNodes(&bindParams.ServiceNodes)

	checkExistDomainErr := services.CheckExistDomain(bindParams.ServiceDomains, []string{})
	if checkExistDomainErr != nil {
		utils.Error(c, checkExistDomainErr.Error())
		return
	}

	checkDomainCertificateErr := services.CheckDomainCertificate(bindParams.Protocol, bindParams.ServiceDomains)
	if checkDomainCertificateErr != nil {
		utils.Error(c, checkDomainCertificateErr.Error())
		return
	}

	createErr := services.ServiceCreate(&bindParams)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceUpdate(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.ServiceAddUpdate{
		IsRelease: utils.IsReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.CorrectServiceTimeOut(&bindParams.Timeouts)
	validators.CorrectServiceDomains(&bindParams.ServiceDomains)
	validators.CorrectServiceAddNodes(&bindParams.ServiceNodes)

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	err := services.CheckExistDomain(bindParams.ServiceDomains, []string{serviceId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkDomainCertificateErr := services.CheckDomainCertificate(bindParams.Protocol, bindParams.ServiceDomains)
	if checkDomainCertificateErr != nil {
		utils.Error(c, checkDomainCertificateErr.Error())
		return
	}

	updateErr := services.ServiceUpdate(serviceId, &bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceInfo(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
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

func ServiceDelete(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkServiceDeleteErr := services.CheckServiceDelete(serviceId)
	if checkServiceDeleteErr != nil {
		utils.Error(c, checkServiceDeleteErr.Error())
		return
	}

	deleteErr := services.ServiceDelete(serviceId)
	if deleteErr != nil {
		utils.Error(c, deleteErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceUpdateName(c *gin.Context) {
	serviceId := c.Param("id")

	var bindParams = validators.ServiceUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	serviceModel := &models.Services{}
	updateErr := serviceModel.ServiceUpdateName(serviceId, bindParams.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchEnable(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.ServiceSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkServiceEnableChangeErr := services.CheckServiceEnableChange(serviceId, bindParams.IsEnable)
	if checkServiceEnableChangeErr != nil {
		utils.Error(c, checkServiceEnableChangeErr.Error())
		return
	}

	serviceModel := &models.Services{}
	updateErr := serviceModel.ServiceSwitchEnable(serviceId, bindParams.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchRelease(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkServiceReleaseErr := services.CheckServiceRelease(serviceId)
	if checkServiceReleaseErr != nil {
		utils.Error(c, checkServiceReleaseErr.Error())
		return
	}

	serviceReleaseErr := services.ServiceRelease(serviceId)
	if serviceReleaseErr != nil {
		utils.Error(c, serviceReleaseErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchWebsocket(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.ServiceSwitchWebsocket{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkServiceWebsocketChangeErr := services.CheckServiceWebsocketChange(serviceId, bindParams.WebSocket)
	if checkServiceWebsocketChangeErr != nil {
		utils.Error(c, checkServiceWebsocketChangeErr.Error())
		return
	}

	serviceModel := &models.Services{}
	updateErr := serviceModel.ServiceSwitchWebsocket(serviceId, bindParams.WebSocket)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchHealthCheck(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.ServiceSwitchHealthCheck{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkServiceHealthCheckChangeErr := services.CheckServiceHealthCheckChange(serviceId, bindParams.HealthCheck)
	if checkServiceHealthCheckChangeErr != nil {
		utils.Error(c, checkServiceHealthCheckChangeErr.Error())
		return
	}

	serviceModel := &models.Services{}
	updateErr := serviceModel.ServiceSwitchHealthCheck(serviceId, bindParams.HealthCheck)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}
