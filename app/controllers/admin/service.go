package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func ServiceAdd(c *gin.Context) {

	var bindParams = &validators.ServiceAddUpdate{
		Release: utils.ReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.CorrectServiceAttributesDefault(bindParams)
	validators.CorrectServiceDomains(bindParams.ServiceDomains)

	s := services.NewServicesService()
	err := s.CheckExistDomain(bindParams.ServiceDomains, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = s.ServiceCreate(bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func ServiceUpdate(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var bindParams = &validators.ServiceAddUpdate{
		Release: utils.ReleaseN,
	}
	if msg, err := packages.ParseRequestParams(c, bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.CorrectServiceDomains(bindParams.ServiceDomains)

	err := services.CheckServiceExist(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	s := services.NewServicesService()

	err = s.CheckExistDomain(bindParams.ServiceDomains, []string{serviceId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := s.ServiceUpdate(serviceId, bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func ServiceInfo(c *gin.Context) {

	serviceId := strings.TrimSpace(c.Param("id"))

	if serviceId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}
	checkServiceExistErr := services.CheckServiceExist(serviceId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	res, err := services.NewServicesService().ServiceInfoById(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	utils.Ok(c, res)
}

func ServiceList(c *gin.Context) {
	var request = &validators.ServiceList{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.NewServicesService().ServiceList(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	res := &utils.ResultPage{
		Param:    request,
		Page:     request.Page,
		PageSize: request.PageSize,
		Data:     list,
		Total:    total,
	}

	utils.Ok(c, res)
}

func ServiceDelete(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	err := services.CheckServiceExist(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = services.NewServicesService().ServiceDelete(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func ServiceUpdateName(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	var request = &validators.ServiceUpdateName{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckServiceExist(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = services.NewServicesService().ServiceUpdateName(serviceId, request)
	if err != nil {
		utils.Error(c, err.Error())
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

	err := services.CheckServiceExist(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = services.NewServicesService().ServiceSwitchEnable(serviceId, bindParams.Enable)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func ServiceSwitchRelease(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("id"))

	err := services.CheckServiceExist(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = services.NewServicesService().ServiceRelease(serviceId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}
