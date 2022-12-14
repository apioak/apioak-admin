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

func ServicePluginConfigList(c *gin.Context) {
	serviceID := strings.TrimSpace(c.Param("service_id"))

	if serviceID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	var request = &validators.ValidatorPluginConfigList{
		Type: models.PluginConfigsTypeService,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	res, err := services.NewPluginsService().PluginConfigList(request.Type, serviceID)

	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	utils.Ok(c, res)
}

func ServicePluginConfigInfo(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("plugin_config_res_id"))

	if pluginConfigID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	res, err := services.NewPluginsService().PluginConfigInfoByResId(pluginConfigID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, res)
}

func ServicePluginConfigAdd(c *gin.Context) {
	var request = &validators.ValidatorPluginConfigAdd{
		Type: models.PluginConfigsTypeService,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}
	pluginConfigResId, err := services.NewPluginsService().PluginConfigAdd(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, map[string]string{
		"res_id": pluginConfigResId,
	})

}

func ServicePluginConfigUpdate(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("plugin_config_res_id"))

	var request = &validators.ValidatorPluginConfigUpdate{
		PluginConfigId: pluginConfigID,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewPluginsService().PluginConfigUpdate(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func ServicePluginConfigSwitchEnable(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("plugin_config_res_id"))

	var request = &validators.ValidatorPluginConfigSwitchEnable{
		PluginConfigId: pluginConfigID,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewPluginsService().PluginConfigSwitchEnable(pluginConfigID, request.Enable)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func ServicePluginConfigDelete(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("plugin_config_res_id"))

	if pluginConfigID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	err := services.NewPluginsService().PluginConfigDelete(pluginConfigID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}
