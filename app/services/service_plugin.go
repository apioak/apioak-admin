package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
	"strings"
)

func PluginCreate(pluginData *validators.ValidatorPluginAdd) error {

	createPluginData := &models.Plugins{
		Name:        pluginData.Name,
		Tag:         pluginData.Tag,
		Icon:        pluginData.Icon,
		Type:        pluginData.Type,
		Description: pluginData.Description,
	}

	pluginModel := models.Plugins{}
	createErr := pluginModel.PluginAdd(createPluginData)

	return createErr
}

func PluginUpdate(id string, pluginUpdate *validators.ValidatorPluginUpdate) error {
	id = strings.TrimSpace(id)

	updatePluginData := models.Plugins{
		Name:        pluginUpdate.Name,
		Icon:        pluginUpdate.Icon,
		Description: pluginUpdate.Description,
	}

	updateErr := updatePluginData.PluginUpdate(id, &updatePluginData)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type StructPluginInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
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
			structPluginInfo.ID = pluginInfo.ID
			structPluginInfo.Name = pluginInfo.Name
			structPluginInfo.Tag = pluginInfo.Tag
			structPluginInfo.Icon = pluginInfo.Icon
			structPluginInfo.Type = pluginInfo.Type
			structPluginInfo.Description = pluginInfo.Description

			pluginInfoList = append(pluginInfoList, structPluginInfo)
		}
	}

	return pluginInfoList, total, pluginInfosErr
}
