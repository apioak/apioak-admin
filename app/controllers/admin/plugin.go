package admin

import (
	"apioak-admin/app/models"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()

	utils.Ok(c, pluginAllTypes)
}

type PluginItem struct {
	ResID       string `json:"res_id"`
	PluginKey   string `json:"plugin_key"`
	Icon        string `json:"icon"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

func PluginAddList(c *gin.Context) {
	pluginModel := models.Plugins{}
	pluginList, err := pluginModel.PluginAllList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	list := make([]PluginItem, 0)
	if len(pluginList) > 0 {
		for _, pluginInfo := range pluginList {
			list = append(list, PluginItem{
				ResID:       pluginInfo.ResID,
				PluginKey:   pluginInfo.PluginKey,
				Icon:        pluginInfo.Icon,
				Type:        pluginInfo.Type,
				Description: pluginInfo.Description,
			})
		}
	}

	utils.Ok(c, list)
}

func PluginInfo(c *gin.Context) {
	pluginResId := strings.TrimSpace(c.Param("plugin_res_id"))

	pluginService := services.PluginsService{}
	pluginConfigDefault, err := pluginService.PluginConfigDefault(pluginResId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, pluginConfigDefault)
}
