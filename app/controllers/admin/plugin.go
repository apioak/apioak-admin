package admin

import (
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"github.com/gin-gonic/gin"
)

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()

	utils.Ok(c, pluginAllTypes)
}

func PluginAddList(c *gin.Context) {
	pluginAddListItem := services.PluginAddListItem{}
	pluginAddList, err := pluginAddListItem.PluginAddList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, pluginAddList)
}
