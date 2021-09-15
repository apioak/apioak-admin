package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
)

func PluginCreate(pluginData *validators.ValidatorPluginAdd) error {

	createPluginData := &models.Plugins{
		Name: pluginData.Name,
		Tag: pluginData.Tag,
		Icon: pluginData.Icon,
		Type: pluginData.Type,
		Description: pluginData.Description,
	}

	pluginModel := models.Plugins{}
	createErr := pluginModel.PluginAdd(createPluginData)

	return createErr
}