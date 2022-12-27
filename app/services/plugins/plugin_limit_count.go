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

func (limitCountConfig PluginLimitCountConfig) PluginConfigParse(configInfo interface{}) (pluginLimitCountConfig interface{}, err error) {
	pluginLimitCount := PluginLimitCount{
		TimeWindow: -999,
		Count:      -999,
	}

	var configInfoJson []byte
	_, ok := configInfo.(string)
	if ok {
		configInfoJson = []byte(fmt.Sprint(configInfo))
	} else {
		configInfoJson, err = json.Marshal(configInfo)
		if err != nil {
			return
		}
	}

	err = json.Unmarshal(configInfoJson, &pluginLimitCount)
	if err != nil {
		return
	}

	pluginLimitCountConfig = pluginLimitCount

	return
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigCheck(configInfo interface{}) error {
	limitCount, err := limitCountConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginLimitCount := limitCount.(PluginLimitCount)

	return limitCountConfig.configValidator(pluginLimitCount)
}

func (limitCountConfig PluginLimitCountConfig) configValidator(config PluginLimitCount) error {

	if config.TimeWindow == -999 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.time_window", "int"))
	}

	if config.Count == -999 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.count", "int"))
	}

	if config.TimeWindow < 1 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.time_window", 1))
	}
	if config.Count < 1 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.count", 1))
	}

	if config.TimeWindow > 86400 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.time_window", 86400))
	}
	if config.Count > 100000000 {
		return errors.New(fmt.Sprintf(
			limitCountValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.count", 100000000))
	}

	return nil
}
