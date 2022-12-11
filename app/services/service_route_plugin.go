package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
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

func CheckRoutePluginDelete(id string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")

	if routePluginInfo.ReleaseStatus == utils.ReleaseStatusY {
		if routePluginInfo.IsEnable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if routePluginInfo.ReleaseStatus == utils.ReleaseStatusT {
		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
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
	if routePluginInfo.ReleaseStatus == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func RoutePluginCreate(routePluginData *validators.RoutePluginAddUpdate) error {
	config, _ := json.Marshal(routePluginData.Config)
	createRoutePlugin := models.RoutePlugins{
		PluginID:      routePluginData.PluginID,
		RouteID:       routePluginData.RouteID,
		Order:         routePluginData.Order,
		Config:        string(config),
		IsEnable:      routePluginData.IsEnable,
		ReleaseStatus: utils.ReleaseStatusU,
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
	routePluginInfo := pluginModel.RoutePluginInfoById(id, "", "")
	if routePluginInfo.ReleaseStatus == utils.ReleaseStatusY {
		createRoutePlugin.ReleaseStatus = utils.ReleaseStatusT
	}

	updateErr := pluginModel.RoutePluginUpdate(id, &createRoutePlugin)

	return updateErr
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
	pluginInfo := pluginModel.PluginInfoByResIdRouteServiceId(routePluginInfo.PluginID)

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(pluginInfo.PluginKey)
	if newPluginContextErr != nil {
		return nil, newPluginContextErr
	}

	parsePluginInfo := newPluginContext.StrategyPluginParse(routePluginInfo.Config)

	return parsePluginInfo, nil
}

func RoutePluginRelease(id string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")
	updateErr := routePluginModel.RoutePluginSwitchRelease(id, utils.ReleaseStatusY)
	if updateErr != nil {
		return updateErr
	}

	configReleaseErr := ServiceRoutePluginConfigRelease(utils.ReleaseTypePush, id)
	if configReleaseErr != nil {
		routePluginModel.RoutePluginSwitchRelease(id, routePluginInfo.ReleaseStatus)
		return configReleaseErr
	}

	return nil
}

func ServiceRoutePluginConfigRelease(releaseType string, id string) error {
	routePluginConfig, routePluginConfigErr := generateRoutePluginConfig(id)
	if routePluginConfigErr != nil {
		return routePluginConfigErr
	}

	routePluginConfigJson, routePluginConfigJsonErr := json.Marshal(routePluginConfig)
	if routePluginConfigJsonErr != nil {
		return routePluginConfigJsonErr
	}
	routePluginConfigStr := string(routePluginConfigJson)

	etcdKey := utils.EtcdKey(utils.EtcdKeyTypePlugin, id)
	if len(etcdKey) == 0 {
		return errors.New(enums.CodeMessages(enums.EtcdKeyNull))
	}

	etcdClient := packages.GetEtcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	defer cancel()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(ctx, etcdKey, routePluginConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(ctx, etcdKey)
	}

	if respErr != nil {
		return errors.New(enums.CodeMessages(enums.EtcdUnavailable))
	}

	return nil
}

type RoutePluginConfig struct {
	ID        string `json:"id"`
	RouteID   string `json:"route_id"`
	PluginTag string `json:"plugin_tag"`
	Order     int    `json:"order"`
	IsEnable  int    `json:"is_enable"`
	Config    string `json:"config"`
}

func generateRoutePluginConfig(id string) (RoutePluginConfig, error) {
	routePluginConfig := RoutePluginConfig{}
	routePluginModel := models.RoutePlugins{}
	routePluginInfo := routePluginModel.RoutePluginInfoById(id, "", "")
	if len(routePluginInfo.ID) == 0 {
		return routePluginConfig, errors.New(enums.CodeMessages(enums.RoutePluginNull))
	}

	pluginModel := models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByResId(routePluginInfo.PluginID)
	if len(pluginInfo.ResID) == 0 {
		return routePluginConfig, errors.New(enums.CodeMessages(enums.PluginNull))
	}

	routePluginConfig.ID = routePluginInfo.ID
	routePluginConfig.RouteID = routePluginInfo.RouteID
	routePluginConfig.PluginTag = pluginInfo.PluginKey
	routePluginConfig.Order = routePluginInfo.Order
	routePluginConfig.IsEnable = routePluginInfo.IsEnable
	routePluginConfig.Config = routePluginInfo.Config

	return routePluginConfig, nil
}
