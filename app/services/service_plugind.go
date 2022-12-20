package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
)

func PluginBasicInfoMaintain() {

	pluginModel := models.Plugins{}
	dbPluginList, _ := pluginModel.PluginAllList()

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
	pluginAllList, err := pluginModel.PluginAllList()
	if len(pluginAllList) == 0 || err != nil {
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
