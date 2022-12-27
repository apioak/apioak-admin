package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type PluginStrategy interface {
	PluginConfigDefault() interface{}
	PluginConfigParse(config interface{}) (interface{}, error)
	PluginConfigCheck(config interface{}) error
}

type PluginContext struct {
	Strategy PluginStrategy
}

func NewPluginContext(pluginTag string) (PluginContext, error) {
	pluginContext := PluginContext{}

	switch strings.ToLower(pluginTag) {
	case utils.PluginKeyCors:
		pluginContext.Strategy = NewCors()
	case utils.PluginKeyKeyAuth:
		pluginContext.Strategy = NewKeyAuth()
	case utils.PluginKeyJwtAuth:
		pluginContext.Strategy = NewJwtAuth()
	case utils.PluginKeyLimitCount:
		pluginContext.Strategy = NewLimitCount()
	case utils.PluginKeyLimitConn:
		pluginContext.Strategy = NewLimitConn()
	case utils.PluginKeyLimitReq:
		pluginContext.Strategy = NewLimitReq()
	default:
		return pluginContext, errors.New(enums.CodeMessages(enums.PluginTagNull))
	}

	return pluginContext, nil
}

func (p PluginContext) StrategyPluginFormatDefault() interface{} {
	return p.Strategy.PluginConfigDefault()
}

func (p PluginContext) StrategyPluginParse(config interface{}) (interface{}, error) {
	return p.Strategy.PluginConfigParse(config)
}

func (p PluginContext) StrategyPluginCheck(config interface{}) error {
	return p.Strategy.PluginConfigCheck(config)
}
