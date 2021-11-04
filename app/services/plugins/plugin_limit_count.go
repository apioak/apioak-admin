package plugins

import (
	"apioak-admin/app/enums"
	"encoding/json"
	"errors"
)

type PluginLimitCountConfig struct {
	ConfigInfo string `json:"config_info"`
}

type PluginLimitCount struct {
	TimeWindow int `json:"time_window"`
	Count      int `json:"count"`
}

func NewLimitCount(config string) PluginLimitCountConfig {
	newLimitCount := PluginLimitCountConfig{
		ConfigInfo: config,
	}

	return newLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginParse() interface{} {
	pluginLimitCount := PluginLimitCount{}
	_ = json.Unmarshal([]byte(limitCountConfig.ConfigInfo), &pluginLimitCount)

	return pluginLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginCheck() error {
	limitCount := limitCountConfig.PluginParse()

	pluginLimitCount := limitCount.(PluginLimitCount)
	if (pluginLimitCount.TimeWindow == 0) || (pluginLimitCount.Count == 0) {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
