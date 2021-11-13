package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
)

func CheckRoutePluginExist(id string, routeId string, pluginId string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, routeId, pluginId)
	if routePluginInfo.ID != id {
		return errors.New(enums.CodeMessages(enums.RoutePluginNull))
	}

	return nil
}

func CheckRoutePluginExistByRoutePluginId(routeId string, pluginId string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoByRoutePluginId(routeId, pluginId)
	if routePluginInfo.RouteID == routeId {
		return errors.New(enums.CodeMessages(enums.RoutePluginExist))
	}

	return nil
}

func CheckRoutePluginEnableOn(id string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")
	if routePluginInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
	}

	return nil
}

func CheckRoutePluginEnableChange(id string, enable int) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")
	if routePluginInfo.IsEnable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func RoutePluginCreate(routePluginData *validators.RoutePluginAddUpdate) error {
	config, _ := json.Marshal(routePluginData.Config)
	createRoutePlugin := models.RoutePlugins{
		PluginID: routePluginData.PluginID,
		RouteID:  routePluginData.RouteID,
		Order:    routePluginData.Order,
		Config:   string(config),
		IsEnable: routePluginData.IsEnable,
	}

	pluginModel := models.RoutePlugins{}
	createErr := pluginModel.RoutePluginAdd(&createRoutePlugin)

	return createErr
}

func RoutePluginUpdate(id string, routePluginData *validators.RoutePluginAddUpdate) error {
	config, _ := json.Marshal(routePluginData.Config)
	createRoutePlugin := models.RoutePlugins{
		Order:    routePluginData.Order,
		Config:   string(config),
		IsEnable: routePluginData.IsEnable,
	}

	pluginModel := models.RoutePlugins{}
	updateErr := pluginModel.RoutePluginUpdate(id, &createRoutePlugin)

	return updateErr
}

func RoutePluginSwitchEnable(id string, enable int) error {
	routePluginModel := models.RoutePlugins{}
	updateErr := routePluginModel.RoutePluginSwitchEnable(id, enable)
	if updateErr != nil {
		return updateErr
	}

	// @todo 触发远程发布数据

	return nil
}

func RoutePluginDelete(id string) error {
	routePluginModel := models.RoutePlugins{}
	deleteErr := routePluginModel.RoutePluginDelete(id)
	if deleteErr != nil {
		return deleteErr
	}

	// @todo 需要同步远程数据中心

	return nil
}

func RoutePluginConfigInfo(id string) (interface{}, error) {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")

	pluginModel := &models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByIdRouteServiceId(routePluginInfo.PluginID)

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(pluginInfo.Tag)
	if newPluginContextErr != nil {
		return nil, newPluginContextErr
	}

	parsePluginInfo := newPluginContext.StrategyPluginParse(routePluginInfo.Config)

	return parsePluginInfo, nil
}
