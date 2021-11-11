package plugins

import (
	"apioak-admin/app/enums"
	"encoding/json"
	"errors"
)

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

func (limitCountConfig PluginLimitCountConfig) PluginConfigParse(configInfo string) interface{} {
	pluginLimitCount := PluginLimitCount{}
	_ = json.Unmarshal([]byte(configInfo), &pluginLimitCount)

	return pluginLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigParseToJson(configInfo string) string {
	limitCount := limitCountConfig.PluginConfigParse(configInfo)
	pluginConfigJson, _ := json.Marshal(limitCount)

	return string(pluginConfigJson)
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigCheck(configInfo string) error {
	limitCount := limitCountConfig.PluginConfigParse(configInfo)
	pluginLimitCount := limitCount.(PluginLimitCount)

	// @todo 增加针对当前插件配置的参数校验

	if (pluginLimitCount.TimeWindow == 0) || (pluginLimitCount.Count == 0) {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
