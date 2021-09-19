package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
	"errors"
)

func CheckRoutePluginExist(routeId string, pluginId string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoByRoutePluginId(routeId, pluginId)
	if routePluginInfo.RouteID == routeId {
		return errors.New(enums.CodeMessages(enums.RoutePluginExist))
	}

	return nil
}

func RoutePluginCreate(routePluginData *validators.RoutePluginAdd) error {

	createRoutePlugin := models.RoutePlugins{
		PluginID: routePluginData.PluginID,
		RouteID:  routePluginData.RouteID,
		Order:    routePluginData.Order,
		Config:   routePluginData.Config,
		IsEnable: routePluginData.IsEnable,
	}

	pluginModel := models.RoutePlugins{}
	createErr := pluginModel.RoutePluginAdd(&createRoutePlugin)

	return createErr
}
