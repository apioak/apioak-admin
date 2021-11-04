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

func (limitCountConfig PluginLimitCountConfig) PluginConfigParse() interface{} {
	pluginLimitCount := PluginLimitCount{}
	_ = json.Unmarshal([]byte(limitCountConfig.ConfigInfo), &pluginLimitCount)

	return pluginLimitCount
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigParseToJson() string {
	limitCount := limitCountConfig.PluginConfigParse()
	pluginConfigJson, _ := json.Marshal(limitCount)

	return string(pluginConfigJson)
}

func (limitCountConfig PluginLimitCountConfig) PluginConfigCheck() error {
	limitCount := limitCountConfig.PluginConfigParse()
	pluginLimitCount := limitCount.(PluginLimitCount)

	// @todo 增加针对当前插件配置的参数校验

	if (pluginLimitCount.TimeWindow == 0) || (pluginLimitCount.Count == 0) {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
