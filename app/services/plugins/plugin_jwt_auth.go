package plugins

import (
	"apioak-admin/app/enums"
	"encoding/json"
	"errors"
)

type PluginJwtAuthConfig struct {
	ConfigInfo string `json:"config_info"`
}

type PluginJwtAuth struct {
	Secret string `json:"secret"`
}

func NewJwtAuth(config string) PluginJwtAuthConfig {
	newJwtAuth := PluginJwtAuthConfig{
		ConfigInfo: config,
	}

	return newJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParse() interface{} {
	pluginJwtAuth := PluginJwtAuth{}
	_ = json.Unmarshal([]byte(jwtAuthConfig.ConfigInfo), &pluginJwtAuth)

	return pluginJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParseToJson() string {
	jwtAuth := jwtAuthConfig.PluginConfigParse()
	pluginConfigJson, _ := json.Marshal(jwtAuth)

	return string(pluginConfigJson)
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigCheck() error {
	jwtAuth := jwtAuthConfig.PluginConfigParse()
	pluginJwtAuth := jwtAuth.(PluginJwtAuth)

	// @todo 增加针对当前插件配置的参数校验

	if len(pluginJwtAuth.Secret) == 0 {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
