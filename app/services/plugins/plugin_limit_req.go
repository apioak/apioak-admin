package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var limitReqValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
		"max":      "[%s] must be %d or less",
		"min":      "[%s] must be %d or greater",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
		"max":      "[%s]必须小于或等于%d",
		"min":      "[%s]最小只能为%d",
	},
}

type PluginLimitReqConfig struct{}

type PluginLimitReq struct {
	Rate  int `json:"rate"`
	Burst int `json:"burst"`
}

func NewLimitReq() PluginLimitReqConfig {
	newLimitReq := PluginLimitReqConfig{}

	return newLimitReq
}

func (limitReqConfig PluginLimitReqConfig) PluginConfigDefault() interface{} {
	pluginLimitReq := PluginLimitReq{
		Rate:             0,
		Burst:            0,
	}

	return pluginLimitReq
}

func (limitReqConfig PluginLimitReqConfig) PluginConfigParse(configInfo interface{}) (pluginLimitReqConfig interface{}, err error) {
	pluginLimitReq := PluginLimitReq{
		Rate:             -999,
		Burst:            -999,
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

	err = json.Unmarshal(configInfoJson, &pluginLimitReq)
	if err != nil {
		return
	}

	pluginLimitReqConfig = pluginLimitReq

	return
}

func (limitReqConfig PluginLimitReqConfig) PluginConfigCheck(configInfo interface{}) error {
	limitReq, _ := limitReqConfig.PluginConfigParse(configInfo)
	pluginLimitReq := limitReq.(PluginLimitReq)

	return limitReqConfig.configValidator(pluginLimitReq)
}

func (limitReqConfig PluginLimitReqConfig) configValidator(config PluginLimitReq) error {

	if config.Rate == -999 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.rate", "int"))
	}
	if config.Burst == -999 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.burst", "int"))
	}

	if config.Rate < 1 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.rate", 1))
	}
	if config.Burst < 0 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.burst", 0))
	}

	if config.Rate > 100000 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.rate", 100000))
	}
	if config.Burst > 5000 {
		return errors.New(fmt.Sprintf(
			limitReqValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.burst", 5000))
	}

	return nil
}
