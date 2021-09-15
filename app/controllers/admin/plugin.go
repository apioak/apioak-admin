package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func PluginAdd(c *gin.Context) {
	var validatorPluginAdd = validators.ValidatorPluginAdd{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginAdd); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetPluginAddAttributesDefault(&validatorPluginAdd)

	pluginModel := models.Plugins{}
	pluginInfos, err := pluginModel.PluginInfosByTags([]string{validatorPluginAdd.Tag}, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	if len(pluginInfos) != 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginTagExist))
		return
	}

	addErr := services.PluginCreate(&validatorPluginAdd)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}
