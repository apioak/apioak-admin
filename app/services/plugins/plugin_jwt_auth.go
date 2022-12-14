package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var jwtAuthValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
		"max":      "[%s] must be %d or less",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
		"max":      "[%s]长度必须小于或等于%d",
	},
}

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

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParse(configInfo interface{}) (interface{}, error) {

	pluginJwtAuth := PluginJwtAuth{
		Secret: "",
	}
	configInfoJson := []byte(fmt.Sprint(configInfo))

	err := json.Unmarshal(configInfoJson, &pluginJwtAuth)

	return pluginJwtAuth, err
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigCheck(configInfo interface{}) error {
	jwtAuth, _ := jwtAuthConfig.PluginConfigParse(configInfo)
	pluginJwtAuth := jwtAuth.(PluginJwtAuth)

	return jwtAuthConfig.configValidator(pluginJwtAuth)
}

func (jwtAuthConfig PluginJwtAuthConfig) configValidator(config PluginJwtAuth) error {

	if len(config.Secret) == 0 {
		return errors.New(fmt.Sprintf(
			jwtAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.secret", "string"))
	}

	if len(config.Secret) > 128 {
		return errors.New(fmt.Sprintf(
			jwtAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.secret", 128))
	}

	return nil
}
