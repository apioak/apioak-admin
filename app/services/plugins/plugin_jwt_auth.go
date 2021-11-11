package plugins

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type PluginJwtAuthConfig struct{}

type PluginJwtAuth struct {
	Secret string `json:"secret"`
}

func NewJwtAuth() PluginJwtAuthConfig {
	newJwtAuth := PluginJwtAuthConfig{}

	return newJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigDefault() interface{} {
	pluginJwtAuth := PluginJwtAuth{
		Secret: utils.Md5(strconv.Itoa(int(time.Now().UnixNano()))),
	}

	return pluginJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParse(configInfo string) interface{} {
	pluginJwtAuth := PluginJwtAuth{}
	_ = json.Unmarshal([]byte(configInfo), &pluginJwtAuth)

	return pluginJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParseToJson(configInfo string) string {
	jwtAuth := jwtAuthConfig.PluginConfigParse(configInfo)
	pluginConfigJson, _ := json.Marshal(jwtAuth)

	return string(pluginConfigJson)
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigCheck(configInfo string) error {
	jwtAuth := jwtAuthConfig.PluginConfigParse(configInfo)
	pluginJwtAuth := jwtAuth.(PluginJwtAuth)

	// @todo 增加针对当前插件配置的参数校验

	if len(pluginJwtAuth.Secret) == 0 {
		return errors.New(enums.CodeMessages(enums.RoutePluginFormatError))
	}

	return nil
}
