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

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()
	
	utils.Ok(c, pluginAllTypes)
}

func PluginAdd(c *gin.Context) {
	var validatorPluginAdd = validators.ValidatorPluginAdd{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginAdd); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetPluginAddAttributesDefault(&validatorPluginAdd)

	pluginModel := models.Plugins{}
	pluginInfos, err := pluginModel.PluginInfosByTags([]string{validatorPluginAdd.Tag}, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	if len(pluginInfos) != 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginTagExist))
		return
	}

	addErr := services.PluginCreate(&validatorPluginAdd)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}

func PluginUpdate(c *gin.Context) {
	pluginId := strings.TrimSpace(c.Param("id"))

	var validatorPluginUpdate = validators.ValidatorPluginUpdate{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginUpdate); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetPluginUpdateAttributesDefault(&validatorPluginUpdate)

	pluginModel := models.Plugins{}
	pluginInfos, pluginInfosErr := pluginModel.PluginInfosByIds([]string{pluginId})
	if pluginInfosErr != nil {
		utils.Error(c, pluginInfosErr.Error())
		return
	}
	if len(pluginInfos) == 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginNull))
		return
	}

	updateErr := services.PluginUpdate(pluginId, &validatorPluginUpdate)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func PluginDelete(c *gin.Context) {
	pluginId := strings.TrimSpace(c.Param("id"))

	pluginModel := models.Plugins{}
	pluginInfos, pluginInfosErr := pluginModel.PluginInfosByIds([]string{pluginId})
	if pluginInfosErr != nil {
		utils.Error(c, pluginInfosErr.Error())
		return
	}

	if len(pluginInfos) == 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginNull))
		return
	}

	routePluginModel := models.RoutePlugins{}
	routePluginInfos, routePluginInfosErr := routePluginModel.RoutePluginInfosByPluginIds([]string{pluginId})
	if routePluginInfosErr != nil {
		utils.Error(c, routePluginInfosErr.Error())
		return
	}
	if len(routePluginInfos) != 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginRouteExist))
		return
	}

	deleteErr := pluginModel.PluginDelete(pluginId)
	if deleteErr != nil {
		utils.Error(c, deleteErr.Error())
		return
	}

	utils.Ok(c)
}

func PluginList(c *gin.Context) {
	var validatorPluginList = validators.PluginList{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginList); err != nil {
		utils.Error(c, msg)
		return
	}

	structPluginInfo := services.StructPluginInfo{}
	routeList, total, err := structPluginInfo.PluginListPage(&validatorPluginList)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = validatorPluginList
	result.Page = validatorPluginList.Page
	result.PageSize = validatorPluginList.PageSize
	result.Total = total
	result.Data = routeList

	utils.Ok(c, result)
}
