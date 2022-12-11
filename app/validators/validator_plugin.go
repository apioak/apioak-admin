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
	pluginTypeTagOneOfErrorMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type ValidatorPluginAdd struct {
	PluginKey   string `json:"plugin_key" zh:"插件标识" en:"Plugin key" binding:"required,min=1,max=30,CheckPluginKeyOneOf"`
	Icon        string `json:"icon" zh:"插件ICON" en:"Plugin icon" binding:"required,min=1,max=30"`
	Type        int    `json:"type" zh:"插件类型" en:"Plugin type" binding:"required,CheckPluginTypeOneOf"`
	Description string `json:"description" zh:"插件描述" en:"Plugin description" binding:"omitempty,max=150"`
}

type ValidatorPluginUpdate struct {
	Icon        string `json:"icon" zh:"插件ICON" en:"Plugin icon" binding:"required,min=1,max=30"`
	Description string `json:"description" zh:"插件描述" en:"Plugin description" binding:"omitempty,max=150"`
}

type PluginList struct {
	Type   int    `form:"type" json:"type" zh:"插件类型" en:"Plugin type" binding:"omitempty,CheckPluginTypeOneOf"`
	Search string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
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
		errMsg = fmt.Sprintf(pluginTypeTagOneOfErrorMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(pluginTypeIds, " "))
		packages.SetAllCustomizeValidatorErrMsgs("CheckPluginTypeOneOf", errMsg)
		return false
	}

	return true
}

func CheckPluginKeyOneOf(fl validator.FieldLevel) bool {
	pluginKey := fl.Field().String()
	pluginAllKeys := utils.PluginAllKeys()

	pluginKeysMap := make(map[string]byte, 0)
	if len(pluginAllKeys) != 0 {
		for _, pluginAllKey := range pluginAllKeys {
			if len(pluginAllKey) == 0 {
				continue
			}

			pluginKeysMap[pluginAllKey] = 0
		}
	}

	_, exist := pluginKeysMap[pluginKey]
	if !exist {
		var errMsg string
		errMsg = fmt.Sprintf(pluginTypeTagOneOfErrorMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(pluginAllKeys, " "))
		packages.SetAllCustomizeValidatorErrMsgs("CheckPluginKeyOneOf", errMsg)
		return false
	}

	return true
}

func GetPluginAddAttributesDefault(pluginAdd *ValidatorPluginAdd) {
	pluginAdd.PluginKey = strings.TrimSpace(pluginAdd.PluginKey)
	pluginAdd.Icon = strings.TrimSpace(pluginAdd.Icon)
	pluginAdd.Description = strings.TrimSpace(pluginAdd.Description)
}

func GetPluginUpdateAttributesDefault(pluginUpdate *ValidatorPluginUpdate) {
	pluginUpdate.Icon = strings.TrimSpace(pluginUpdate.Icon)
	pluginUpdate.Description = strings.TrimSpace(pluginUpdate.Description)
}
