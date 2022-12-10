package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/validators"
	"errors"
	"strings"
)

func CheckPluginExist(pluginResId string) error {
	pluginModel := &models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByResIdRouteServiceId(pluginResId)
	if pluginInfo.ResID != pluginResId {
		return errors.New(enums.CodeMessages(enums.PluginNull))
	}

	return nil
}

func CheckPluginConfig(pluginId string, pluginConfig *validators.RoutePluginAddUpdate) error {
	pluginModel := &models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByResIdRouteServiceId(pluginId)

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(pluginInfo.Key)
	if newPluginContextErr != nil {
		return newPluginContextErr
	}

	pluginCheckErr := newPluginContext.StrategyPluginCheck(pluginConfig.Config)
	if pluginCheckErr != nil {
		return pluginCheckErr
	}

	pluginConfig.Config = newPluginContext.StrategyPluginParse(pluginConfig.Config)

	return nil
}

func PluginCreate(pluginData *validators.ValidatorPluginAdd) error {

	createPluginData := &models.Plugins{
		Name:        pluginData.Name,
		Key:         pluginData.Key,
		Icon:        pluginData.Icon,
		Type:        pluginData.Type,
		Description: pluginData.Description,
	}

	pluginModel := models.Plugins{}
	createErr := pluginModel.PluginAdd(createPluginData)

	return createErr
}

func PluginUpdate(resId string, pluginUpdate *validators.ValidatorPluginUpdate) error {
	resId = strings.TrimSpace(resId)

	updatePluginData := models.Plugins{
		Name:        pluginUpdate.Name,
		Icon:        pluginUpdate.Icon,
		Description: pluginUpdate.Description,
	}

	updateErr := updatePluginData.PluginUpdate(resId, &updatePluginData)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type StructPluginInfo struct {
	ResID       string `json:"res_id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Icon        string `json:"icon"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

func (s *StructPluginInfo) PluginListPage(param *validators.PluginList) ([]StructPluginInfo, int, error) {

	pluginModel := models.Plugins{}
	pluginInfos, total, pluginInfosErr := pluginModel.PluginListPage(param)

	pluginInfoList := make([]StructPluginInfo, 0)
	if len(pluginInfos) != 0 {
		for _, pluginInfo := range pluginInfos {
			structPluginInfo := StructPluginInfo{}
			structPluginInfo.ResID = pluginInfo.ResID
			structPluginInfo.Name = pluginInfo.Name
			structPluginInfo.Key = pluginInfo.Key
			structPluginInfo.Icon = pluginInfo.Icon
			structPluginInfo.Type = pluginInfo.Type
			structPluginInfo.Description = pluginInfo.Description

			pluginInfoList = append(pluginInfoList, structPluginInfo)
		}
	}

	return pluginInfoList, total, pluginInfosErr
}

type PluginInfoService struct {
	ResID       string 		`json:"res_id"`
	Name        string      `json:"name"`
	Key         string      `json:"key"`
	Icon        string      `json:"icon"`
	Type        int         `json:"type"`
	Description string      `json:"description"`
	Config      interface{} `json:"config"`
}

func (p *PluginInfoService) PluginInfoByResId(resId string) (PluginInfoService, error) {
	pluginInfo := PluginInfoService{}

	pluginModel := models.Plugins{}
	plugin := pluginModel.PluginInfoByResId(resId)
	if len(plugin.ResID) == 0 {
		return pluginInfo, errors.New(enums.CodeMessages(enums.PluginNull))
	}

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(plugin.Key)
	if newPluginContextErr != nil {
		return pluginInfo, newPluginContextErr
	}

	pluginConfig := newPluginContext.StrategyPluginFormatDefault()

	pluginInfo.ResID = plugin.ResID
	pluginInfo.Name = plugin.Name
	pluginInfo.Key = plugin.Key
	pluginInfo.Icon = plugin.Icon
	pluginInfo.Type = plugin.Type
	pluginInfo.Description = plugin.Description
	pluginInfo.Config = pluginConfig

	return pluginInfo, nil
}
