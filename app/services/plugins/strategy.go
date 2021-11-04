package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type PluginStrategy interface {
	PluginConfigParse() interface{}
	PluginConfigParseToJson() string
	PluginConfigCheck() error
}

type PluginContext struct {
	Strategy PluginStrategy
}

func NewPluginContext(pluginTag string, config string) (PluginContext, error) {
	pluginContext := PluginContext{}

	switch strings.ToLower(pluginTag) {
	case utils.PluginTagNameJwtAuth:
		pluginContext.Strategy = NewJwtAuth(config)
	case utils.PluginTagNameLimitCount:
		pluginContext.Strategy = NewLimitCount(config)
	default:
		return pluginContext, errors.New(enums.CodeMessages(enums.PluginTagNull))
	}

	return pluginContext, nil
}

func (p PluginContext) StrategyPluginParse() interface{} {
	return p.Strategy.PluginConfigParse()
}

func (p PluginContext) StrategyPluginJson() string {
	return p.Strategy.PluginConfigParseToJson()
}

func (p PluginContext) StrategyPluginCheck() error {
	return p.Strategy.PluginConfigCheck()
}
