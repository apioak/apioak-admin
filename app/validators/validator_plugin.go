package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

var (
	pluginTypeOneOfErrorMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type ValidatorPluginAdd struct {
	Name        string `json:"name" zh:"插件名称" en:"Plugin name" binding:"required,min=1,max=30"`
	Tag         string `json:"tag" zh:"插件标识" en:"Plugin tag" binding:"required,min=1,max=30"`
	Icon        string `json:"icon" zh:"插件ICON" en:"Plugin icon" binding:"required,min=1,max=30"`
	Type        int    `json:"type" zh:"插件类型" en:"Plugin type" binding:"required,CheckPluginTypeOneOf"`
	Description string `json:"description" zh:"插件描述" en:"Plugin description" binding:"omitempty,max=150"`
}

func CheckPluginTypeOneOf(fl validator.FieldLevel) bool {
	pluginTypeId := fl.Field().Int()

	pluginAllTypes := utils.PluginAllTypes()

	pluginTypeIdsMap := make(map[int]byte, 0)
	pluginTypeIds := make([]string, 0)
	if len(pluginAllTypes) != 0 {
		for _, pluginAllType := range pluginAllTypes {
			if pluginAllType.Id == 0 {
				continue
			}

			pluginTypeIds = append(pluginTypeIds, strconv.Itoa(pluginAllType.Id))
			pluginTypeIdsMap[pluginAllType.Id] = 0
		}
	}

	_, exist := pluginTypeIdsMap[int(pluginTypeId)]
	if !exist {
		var errMsg string
		errMsg = fmt.Sprintf(pluginTypeOneOfErrorMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(pluginTypeIds, " "))
		packages.SetAllCustomizeValidatorErrMsgs("CheckPluginTypeOneOf", errMsg)
		return false
	}

	return true
}

func GetPluginAddAttributesDefault(pluginAdd *ValidatorPluginAdd) {
	pluginAdd.Name = strings.TrimSpace(pluginAdd.Name)
	pluginAdd.Tag = strings.TrimSpace(pluginAdd.Tag)
	pluginAdd.Icon = strings.TrimSpace(pluginAdd.Icon)
	pluginAdd.Description = strings.TrimSpace(pluginAdd.Description)
}
