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

type ValidatorPluginConfigAdd struct {
	Name     string      `json:"name" zh:"插件名称" en:"Plugin name" binding:"omitempty,min=1,max=30"`
	PluginID string      `json:"plugin_id" zh:"插件ID" en:"Plugin ID" binding:"required"`
	Type     int         `json:"type" zh:"资源类型" en:"Resource type" binding:"omitempty,oneof=1 2"`
	TargetID string      `json:"target_id" zh:"资源ID" en:"Resource ID" binding:"required"`
	Enable   int         `json:"enable" zh:"插件开关" en:"Plugin enable" binding:"omitempty"`
	Config   interface{} `json:"config" zh:"插件配置" en:"Plugin config" binding:"omitempty"`
}

type ValidatorPluginConfigUpdate struct {
	PluginConfigId string      `json:"plugin_config_id" zh:"插件配置ID" en:"Plugin config ID" binding:"required"`
	Name           string      `json:"name" zh:"插件名称" en:"Plugin name" binding:"omitempty,min=1,max=30"`
	Config         interface{} `json:"config" zh:"插件配置" en:"Plugin config" binding:"omitempty"`
}

type ValidatorPluginConfigSwitchEnable struct {
	PluginConfigId string `json:"plugin_config_id" zh:"插件配置ID" en:"Plugin config ID" binding:"required"`
	Enable         int    `json:"enable" zh:"插件开关" en:"plugin enable" binding:"required,oneof=1 2"`
}

type ValidatorPluginConfigList struct {
	Type int `form:"type" json:"type" zh:"资源类型" en:"Resource type" binding:"omitempty,oneof=1 2"`
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
