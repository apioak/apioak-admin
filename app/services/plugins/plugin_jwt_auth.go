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

func (jwtAuthConfig PluginJwtAuthConfig) PluginParse() interface{} {
	pluginJwtAuth := PluginJwtAuth{}
	_ = json.Unmarshal([]byte(jwtAuthConfig.ConfigInfo), &pluginJwtAuth)

	return pluginJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginCheck() error {
	jwtAuth := jwtAuthConfig.PluginParse()

	pluginJwtAuth := jwtAuth.(PluginJwtAuth)
	if len(pluginJwtAuth.Secret) == 0 {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
