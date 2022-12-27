package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var limitConnValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
		"max":      "[%s] must be %d or less",
		"min":      "[%s] must be %d or greater",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
		"max":      "[%s]必须小于或等于%d",
		"min":      "[%s]必须大于或等于%d",
	},
}

type PluginLimitConnConfig struct{}

type PluginLimitConn struct {
	Rate             int `json:"rate"`
	Burst            int `json:"burst"`
	DefaultConnDelay int `json:"default_conn_delay"`
}

func NewLimitConn() PluginLimitConnConfig {
	newLimitConn := PluginLimitConnConfig{}

	return newLimitConn
}

func (limitConnConfig PluginLimitConnConfig) PluginConfigDefault() interface{} {
	pluginLimitConn := PluginLimitConn{
		Rate:             0,
		Burst:            0,
		DefaultConnDelay: 0,
	}

	return pluginLimitConn
}

func (limitConnConfig PluginLimitConnConfig) PluginConfigParse(configInfo interface{}) (pluginLimitConnConfig interface{}, err error) {
	pluginLimitConn := PluginLimitConn{
		Rate:             -999,
		Burst:            -999,
		DefaultConnDelay: -999,
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

	err = json.Unmarshal(configInfoJson, &pluginLimitConn)
	if err != nil {
		return
	}
	
	pluginLimitConnConfig = pluginLimitConn

	return
}

func (limitConnConfig PluginLimitConnConfig) PluginConfigCheck(configInfo interface{}) error {
	limitConn, err := limitConnConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginLimitConn := limitConn.(PluginLimitConn)

	return limitConnConfig.configValidator(pluginLimitConn)
}

func (limitConnConfig PluginLimitConnConfig) configValidator(config PluginLimitConn) error {

	if config.Rate == -999 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.rate", "int"))
	}
	if config.Burst == -999 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.burst", "int"))
	}
	if config.DefaultConnDelay == -999 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.default_conn_delay", "int"))
	}

	if config.Rate < 1 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.rate", 1))
	}
	if config.Burst < 1 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.burst", 1))
	}
	if config.DefaultConnDelay < 1 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.default_conn_delay", 1))
	}

	if config.Rate > 100000 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.rate", 100000))
	}
	if config.Burst > 50000 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.burst", 50000))
	}
	if config.DefaultConnDelay > 60 {
		return errors.New(fmt.Sprintf(
			limitConnValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.default_conn_delay", 60))
	}

	return nil
}
