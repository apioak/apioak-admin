package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type PluginStrategy interface {
	PluginParse() interface{}
	PluginCheck() error
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
	return p.Strategy.PluginParse()
}

func (p PluginContext) StrategyPluginCheck() error {
	return p.Strategy.PluginCheck()
}
