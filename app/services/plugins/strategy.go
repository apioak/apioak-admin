package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type PluginStrategy interface {
	PluginConfigDefault() interface{}
	PluginConfigParse(config interface{}) interface{}
	PluginConfigCheck(config interface{}) error
}

type PluginContext struct {
	Strategy PluginStrategy
}

func NewPluginContext(pluginTag string) (PluginContext, error) {
	pluginContext := PluginContext{}

	switch strings.ToLower(pluginTag) {
	case utils.PluginKeyNameJwtAuth:
		pluginContext.Strategy = NewJwtAuth()
	case utils.PluginKeyNameLimitCount:
		pluginContext.Strategy = NewLimitCount()
	default:
		return pluginContext, errors.New(enums.CodeMessages(enums.PluginTagNull))
	}

	return pluginContext, nil
}

func (p PluginContext) StrategyPluginFormatDefault() interface{} {
	return p.Strategy.PluginConfigDefault()
}

func (p PluginContext) StrategyPluginParse(config interface{}) interface{} {
	return p.Strategy.PluginConfigParse(config)
}

func (p PluginContext) StrategyPluginCheck(config interface{}) error {
	return p.Strategy.PluginConfigCheck(config)
}
