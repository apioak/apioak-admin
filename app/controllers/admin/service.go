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

// @todo 服务列表
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
	result.Param= serviceListStruct
	result.Page = serviceListStruct.Page
	result.PageSize = serviceListStruct.PageSize
	result.Total = total
	result.Data = serviceList

	utils.Ok(c, result)
}

// @todo 服务详情
func ServiceInfo(c *gin.Context) {

}

func ServiceAdd(c *gin.Context) {

	var serviceAddUpdateStruct = validators.ServiceAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &serviceAddUpdateStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckExistDomain(serviceAddUpdateStruct.ServiceNames, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceAddUpdateStruct.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateStruct.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateStruct.ServiceNames)
	serviceNodes := validators.GetServiceAddNodes(serviceAddUpdateStruct.ServiceNodes)

	createErr := services.ServiceCreate(&serviceAddUpdateStruct, &serviceDomains, &serviceNodes)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
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

	err := services.CheckExistDomain(serviceAddUpdateStruct.ServiceNames, []string{serviceId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceAddUpdateStruct.Timeouts = validators.GetServiceAddTimeOut(serviceAddUpdateStruct.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceAddUpdateStruct.ServiceNames)
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
