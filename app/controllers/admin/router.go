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

func RouterAdd(c *gin.Context) {
	bindParams := validators.ValidatorRouterAddUpdate{
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.GetRouterAttributesDefault(&bindParams)
	validators.CorrectUpstreamDefault(&bindParams.UpstreamAddUpdate)
	validators.CorrectUpstreamAddNodes(&bindParams.UpstreamNodes)

	if bindParams.RouterPath == utils.DefaultRouterPath {
		utils.Error(c, enums.CodeMessages(enums.RouterDefaultPathNoPermission))
		return
	}

	checkServiceExistErr := services.CheckServiceExist(bindParams.ServiceResID)
	if (len(bindParams.ServiceResID) == 0) || (checkServiceExistErr != nil) {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	err := services.CheckExistServiceRouterPath(bindParams.ServiceResID, bindParams.RouterPath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	_, createErr := services.RouterCreate(&bindParams)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterList(c *gin.Context) {
	var bindParams = validators.ValidatorRouterList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(bindParams.ServiceResID)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	structRouterList := services.RouterListItem{}
	routerList, total, err := structRouterList.RouterListPage(bindParams.ServiceResID, &bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = bindParams
	result.Page = bindParams.Page
	result.PageSize = bindParams.PageSize
	result.Total = total
	result.Data = routerList

	utils.Ok(c, result)
}

func RouterInfo(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	structRouterInfo := services.StructRouterInfo{}
	routeInfo, routerInfoErr := structRouterInfo.RouterInfoByServiceRouterId(serviceResId, routerResId)
	if routerInfoErr != nil {
		utils.Error(c, enums.CodeMessages(enums.RouterNull))
		return
	}

	utils.Ok(c, routeInfo)
}

func RouterUpdate(c *gin.Context) {
	var bindParams = validators.ValidatorRouterAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouterAttributesDefault(&bindParams)
	validators.CorrectUpstreamDefault(&bindParams.UpstreamAddUpdate)
	validators.CorrectUpstreamAddNodes(&bindParams.UpstreamNodes)

	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	// @todo 默认服务下的全量路由检测，后续版本中放开
	// checkEditDefaultPathRouterErr := services.CheckEditDefaultPathRouter(routerResId)
	// if checkEditDefaultPathRouterErr != nil {
	// 	utils.Error(c, checkEditDefaultPathRouterErr.Error())
	// 	return
	// }

	checkServiceRouterPathErr := services.CheckServiceRouterPath(bindParams.RouterPath)
	if checkServiceRouterPathErr != nil {
		utils.Error(c, checkServiceRouterPathErr.Error())
		return
	}

	err := services.CheckExistServiceRouterPath(bindParams.ServiceResID, bindParams.RouterPath, []string{serviceResId})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := services.RouterUpdate(routerResId, bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterUpdateName(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	var bindParams = validators.RouterUpdateName{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	routerModel := models.Routers{}
	updateErr := routerModel.RouterUpdateName(routerResId, bindParams.Name)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterSwitchEnable(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	var bindParams = validators.RouterSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRouter(routerResId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	checkRouteEnableChangeErr := services.CheckRouterEnableChange(routerResId, bindParams.Enable)
	if checkRouteEnableChangeErr != nil {
		utils.Error(c, checkRouteEnableChangeErr.Error())
		return
	}

	routerModel := models.Routers{}
	updateErr := routerModel.RouterSwitchEnable(routerResId, bindParams.Enable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterDelete(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routeResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouteErr := services.CheckRouterExist(routeResId, serviceResId)
	if checkExistRouteErr != nil {
		utils.Error(c, checkExistRouteErr.Error())
		return
	}

	checkEditDefaultPathRouteErr := services.CheckEditDefaultPathRouter(routeResId)
	if checkEditDefaultPathRouteErr != nil {
		utils.Error(c, checkEditDefaultPathRouteErr.Error())
		return
	}

	deleteErr := services.RouterDelete(routeResId)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func RouterSwitchRelease(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	// @todo 这里增加检测服务是否已发布，只要不是未发布都可以通过，请优先发布服务（未发布的服务不允许发布路由）

	checkExistRouterErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouterErr != nil {
		utils.Error(c, checkExistRouterErr.Error())
		return
	}

	checkRouterReleaseErr := services.CheckRouterRelease(routerResId)
	if checkRouterReleaseErr != nil {
		utils.Error(c, checkRouterReleaseErr.Error())
		return
	}

	serviceRouterReleaseErr := services.RouterUpstreamRelease([]string{routerResId}, utils.ReleaseTypePush)
	if serviceRouterReleaseErr != nil {
		utils.Error(c, serviceRouterReleaseErr.Error())
		return
	}

	utils.Ok(c)
}

func RouterCopy(c *gin.Context) {
	serviceResId := strings.TrimSpace(c.Param("service_res_id"))
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	checkServiceExistErr := services.CheckServiceExist(serviceResId)
	if checkServiceExistErr != nil {
		utils.Error(c, checkServiceExistErr.Error())
		return
	}

	checkExistRouterErr := services.CheckRouterExist(routerResId, serviceResId)
	if checkExistRouterErr != nil {
		utils.Error(c, checkExistRouterErr.Error())
		return
	}

	err := services.RouterCopy(routerResId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func RouterPluginConfigAdd(c *gin.Context) {
	var request = &validators.ValidatorPluginConfigAdd{
		Type: models.PluginConfigsTypeRouter,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	_, err := services.NewPluginsService().PluginConfigAdd(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func RouterPluginConfigList(c *gin.Context) {
	routerResId := strings.TrimSpace(c.Param("router_res_id"))

	if routerResId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	var request = &validators.ValidatorPluginConfigList{
		Type: models.PluginConfigsTypeRouter,
	}

	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	res, err := services.NewPluginsService().PluginConfigList(request.Type, routerResId)

	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, res)
}

func RouterPluginConfigInfo(c *gin.Context) {
	pluginConfigResId := strings.TrimSpace(c.Param("res_id"))

	if pluginConfigResId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	res, err := services.NewPluginsService().PluginConfigInfoByResId(pluginConfigResId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, res)
}

func RouterPluginConfigUpdate(c *gin.Context) {
	pluginConfigResId := strings.TrimSpace(c.Param("res_id"))

	var request = &validators.ValidatorPluginConfigUpdate{
		PluginConfigId: pluginConfigResId,
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

func RouterPluginConfigDelete(c *gin.Context) {
	pluginConfigResId := strings.TrimSpace(c.Param("res_id"))

	if pluginConfigResId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	err := services.NewPluginsService().PluginConfigDelete(pluginConfigResId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func RouterPluginConfigSwitchEnable(c *gin.Context) {
	pluginConfigResId := strings.TrimSpace(c.Param("res_id"))

	var request = &validators.ValidatorPluginConfigSwitchEnable{
		PluginConfigId: pluginConfigResId,
	}

	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewPluginsService().PluginConfigSwitchEnable(pluginConfigResId, request.Enable)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}


