package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var limitCountValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
		"max":      "[%s] must be %d or less",
		"min":      "[%s] must be %d or greater",
		"type":     "[%s] Type error, expected type: %s",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
		"max":      "[%s]必须小于或等于%d",
		"min":      "[%s]最小只能为%d",
		"type":     "[%s]类型错误，期望类型:%s",
	},
}

type PluginLimitCountConfig struct{}

type PluginLimitCount struct {
	TimeWindow int `json:"time_window"`
	Count      int `json:"count"`
}

func NewLimitCount() PluginLimitCountConfig {
	newLimitCount := PluginLimitCountConfig{}

	return newLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigDefault() interface{} {
	pluginLimitCount := PluginLimitCount{
		TimeWindow: 60,
		Count:      1000,
	}

	return pluginLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigParse(configInfo interface{}) (interface{}, error) {
	pluginLimitCount := PluginLimitCount{
		TimeWindow: -9999999,
		Count:      -9999999,
	}

	configInfoJson := []byte(fmt.Sprint(configInfo))

	err := json.Unmarshal(configInfoJson, &pluginLimitCount)

	return pluginLimitCount, err
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigCheck(configInfo interface{}) error {
	limitCount, _ := limitCountConfig.PluginConfigParse(configInfo)
	pluginLimitCount := limitCount.(PluginLimitCount)

	return limitCountConfig.configValidator(pluginLimitCount)
}

func (limitCountConfig PluginLimitCountConfig) configValidator(config PluginLimitCount) error {

	if config.TimeWindow == -9999999 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.time_window", "int"))
	}
	if config.TimeWindow < 0 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.time_window", 0))
	}

	if config.Count == -9999999 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.count", "int"))
	}
	if config.Count < 0 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.count", 0))
	}

	return nil
}
