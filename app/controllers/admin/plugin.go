package admin

import (
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"github.com/gin-gonic/gin"
)

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()

	utils.Ok(c, pluginAllTypes)
}

func PluginAddList(c *gin.Context) {
	pluginModel := models.Plugins{}
	pluginList, err := pluginModel.PluginAllList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, pluginList)
}
