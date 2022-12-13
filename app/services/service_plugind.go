package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strings"
)

func CheckPluginExist(pluginResId string) error {
	pluginModel := &models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByResIdRouterServiceId(pluginResId)
	if pluginInfo.ResID != pluginResId {
		return errors.New(enums.CodeMessages(enums.PluginNull))
	}

	return nil
}

func CheckPluginConfig(pluginId string, pluginConfig *validators.RouterPluginAddUpdate) error {
	pluginModel := &models.Plugins{}
	pluginInfo := pluginModel.PluginInfoByResIdRouterServiceId(pluginId)

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(pluginInfo.PluginKey)
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

func PluginUpdate(resId string, pluginUpdate *validators.ValidatorPluginUpdate) error {
	resId = strings.TrimSpace(resId)

	updatePluginData := models.Plugins{
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
			structPluginInfo.Key = pluginInfo.PluginKey
			structPluginInfo.Icon = pluginInfo.Icon
			structPluginInfo.Type = pluginInfo.Type
			structPluginInfo.Description = pluginInfo.Description

			pluginInfoList = append(pluginInfoList, structPluginInfo)
		}
	}

	return pluginInfoList, total, pluginInfosErr
}

type PluginInfoService struct {
	ResID       string      `json:"res_id"`
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

	newPluginContext, newPluginContextErr := plugins.NewPluginContext(plugin.PluginKey)
	if newPluginContextErr != nil {
		return pluginInfo, newPluginContextErr
	}

	pluginConfig := newPluginContext.StrategyPluginFormatDefault()

	pluginInfo.ResID = plugin.ResID
	pluginInfo.Key = plugin.PluginKey
	pluginInfo.Icon = plugin.Icon
	pluginInfo.Type = plugin.Type
	pluginInfo.Description = plugin.Description
	pluginInfo.Config = pluginConfig

	return pluginInfo, nil
}

func PluginBasicInfoMaintain() {

	pluginModel := models.Plugins{}
	dbPluginList := pluginModel.PluginAllList()

	dbPluginMapResId := make(map[string]models.Plugins)

	for _, dbPluginInfo := range dbPluginList {
		dbPluginMapResId[dbPluginInfo.ResID] = dbPluginInfo
	}

	configPluginList := utils.AllConfigPluginData()

	for _, configPluginInfo := range configPluginList {

		dbPluginMapInfo, ok := dbPluginMapResId[configPluginInfo.ResID]

		if ok {
			if (configPluginInfo.PluginKey != dbPluginMapInfo.PluginKey) ||
				(configPluginInfo.Type != dbPluginMapInfo.Type) {

				dbPluginMapInfo.PluginKey = configPluginInfo.PluginKey
				dbPluginMapInfo.Type = configPluginInfo.Type
				pluginModel.PluginUpdate(configPluginInfo.ResID, &dbPluginMapInfo)
			}

		} else {

			pluginModel.PluginDelByPluginKeys([]string{configPluginInfo.PluginKey}, []string{})

			newPluginData := pluginModel
			newPluginData.ResID = configPluginInfo.ResID
			newPluginData.Type = configPluginInfo.Type
			newPluginData.PluginKey = configPluginInfo.PluginKey
			newPluginData.Icon = configPluginInfo.Icon

			pluginModel.PluginAdd(&newPluginData)
		}
	}
}

type PluginAddListItem struct {
	ResID       string `json:"res_id"`
	PluginKey   string `json:"plugin_key"`
	Icon        string `json:"icon"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

func (s PluginAddListItem) PluginAddList() (list []PluginAddListItem, err error) {
	pluginModel := models.Plugins{}
	pluginAllList := pluginModel.PluginAllList()
	if len(pluginAllList) == 0 {
		return
	}

	for _, pluginInfo := range pluginAllList {
		list = append(list, PluginAddListItem{
			ResID:       pluginInfo.ResID,
			PluginKey:   pluginInfo.PluginKey,
			Icon:        pluginInfo.Icon,
			Type:        pluginInfo.Type,
			Description: pluginInfo.Description,
		})
	}

	return
}
