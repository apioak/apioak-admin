package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
	"fmt"
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

func CheckRoutePluginDelete(id string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")

	if routePluginInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
	}

	if routePluginInfo.IsRelease != utils.IsReleaseY {
		return errors.New(enums.CodeMessages(enums.EnablePublishedONOp))
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

func CheckRoutePluginRelease(id string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")
	if routePluginInfo.IsRelease == utils.IsReleaseY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func RoutePluginCreate(routePluginData *validators.RoutePluginAddUpdate) error {
	config, _ := json.Marshal(routePluginData.Config)
	createRoutePlugin := models.RoutePlugins{
		PluginID:  routePluginData.PluginID,
		RouteID:   routePluginData.RouteID,
		Order:     routePluginData.Order,
		Config:    string(config),
		IsEnable:  routePluginData.IsEnable,
		IsRelease: utils.IsReleaseN,
	}

	pluginModel := models.RoutePlugins{}
	createErr := pluginModel.RoutePluginAdd(&createRoutePlugin)

	return createErr
}

func RoutePluginUpdate(id string, routePluginData *validators.RoutePluginAddUpdate) error {
	config, _ := json.Marshal(routePluginData.Config)
	createRoutePlugin := models.RoutePlugins{
		Order:     routePluginData.Order,
		Config:    string(config),
		IsEnable:  routePluginData.IsEnable,
		IsRelease: utils.IsReleaseN,
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

	return nil
}

func RoutePluginDelete(id string) error {
	configReleaseErr := ServiceRoutePluginConfigRelease(utils.ReleaseTypeDelete, id)
	if configReleaseErr != nil {
		return configReleaseErr
	}

	routePluginModel := models.RoutePlugins{}
	deleteErr := routePluginModel.RoutePluginDelete(id)
	if deleteErr != nil {
		ServiceRoutePluginConfigRelease(utils.ReleaseTypePush, id)
		return deleteErr
	}

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

func RoutePluginRelease(id string) error {
	routePluginModel := models.RoutePlugins{}
	updateErr := routePluginModel.RoutePluginSwitchRelease(id, utils.IsReleaseY)
	if updateErr != nil {
		return updateErr
	}

	configReleaseErr := ServiceRoutePluginConfigRelease(utils.ReleaseTypePush, id)
	if configReleaseErr != nil {
		routePluginModel.RoutePluginSwitchRelease(id, utils.IsReleaseN)
		return configReleaseErr
	}

	return nil
}

func ServiceRoutePluginConfigRelease(releaseType string, id string) error {

	// @todo 获取指定服务路由插件的配置数据
	//serviceRouteConfig := generateServicesRoutePluginConfig(serviceId)

	// @todo 获取数据注册中心对应 服务配置 的key

	fmt.Println("=========service route plugin release:", releaseType, id)

	// @todo 发布配置到 数据注册中心

	return nil
}

func generateServicesRoutePluginConfig(id string) string {

	// @todo 根据服务ID 拼接服务的配置数据（主要是用于同步到数据面使用）

	return ""
}
