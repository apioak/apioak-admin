package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type PluginStrategy interface {
	PluginConfigDefault() interface{}
	PluginConfigParse(config string) interface{}
	PluginConfigParseToJson(config string) string
	PluginConfigCheck(config string) error
}

type PluginContext struct {
	Strategy PluginStrategy
}

func NewPluginContext(pluginTag string) (PluginContext, error) {
	pluginContext := PluginContext{}

	switch strings.ToLower(pluginTag) {
	case utils.PluginTagNameJwtAuth:
		pluginContext.Strategy = NewJwtAuth()
	case utils.PluginTagNameLimitCount:
		pluginContext.Strategy = NewLimitCount()
	default:
		return pluginContext, errors.New(enums.CodeMessages(enums.PluginTagNull))
	}

	return pluginContext, nil
}

func (p PluginContext) StrategyPluginFormatDefault() interface{} {
	return p.Strategy.PluginConfigDefault()
}

func (p PluginContext) StrategyPluginParse(config string) interface{} {
	return p.Strategy.PluginConfigParse(config)
}

func (p PluginContext) StrategyPluginJson(config string) string {
	return p.Strategy.PluginConfigParseToJson(config)
}

func (p PluginContext) StrategyPluginCheck(config string) error {
	return p.Strategy.PluginConfigCheck(config)
}
