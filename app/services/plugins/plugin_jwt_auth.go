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
	JwtKey string `json:"jwt_key"`
}

func NewJwtAuth() PluginJwtAuthConfig {
	newJwtAuth := PluginJwtAuthConfig{}

	return newJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigDefault() interface{} {
	pluginJwtAuth := PluginJwtAuth{
		JwtKey: utils.Md5(strconv.Itoa(int(time.Now().UnixNano()))),
	}

	return pluginJwtAuth
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigParse(configInfo interface{}) (pluginJwtAuthConfig interface{}, err error) {

	pluginJwtAuth := PluginJwtAuth{
		JwtKey: "",
	}

	var configInfoJson []byte
	_, ok := configInfo.(string)
	if ok {
		configInfoJson = []byte(fmt.Sprint(configInfo))
	} else {
		configInfoJson, err = json.Marshal(configInfo)
		if err != nil {
			return
		}
	}

	err = json.Unmarshal(configInfoJson, &pluginJwtAuth)
	if err != nil {
		return
	}

	pluginJwtAuthConfig = pluginJwtAuth

	return
}

func (jwtAuthConfig PluginJwtAuthConfig) PluginConfigCheck(configInfo interface{}) error {
	jwtAuth, _ := jwtAuthConfig.PluginConfigParse(configInfo)
	pluginJwtAuth := jwtAuth.(PluginJwtAuth)

	return jwtAuthConfig.configValidator(pluginJwtAuth)
}

func (jwtAuthConfig PluginJwtAuthConfig) configValidator(config PluginJwtAuth) error {

	if len(config.JwtKey) == 0 {
		return errors.New(fmt.Sprintf(
			jwtAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.jwt_key", "string"))
	}

	if len(config.JwtKey) < 10 {
		return errors.New(fmt.Sprintf(
			jwtAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.jwt_key", 10))
	}

	if len(config.JwtKey) > 32 {
		return errors.New(fmt.Sprintf(
			jwtAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.jwt_key", 32))
	}

	return nil
}
