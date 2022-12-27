package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var keyAuthValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
		"max":      "[%s] Length must be less than or equal to %d",
		"min":      "[%s] Length must be less than or equal to %d",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
		"max":      "[%s]长度必须小于或等于%d",
		"min":      "[%s]长度必须大于或等于%d",
	},
}

type PluginKeyAuthConfig struct{}

type PluginKeyAuth struct {
	Secret string `json:"secret"`
}

func NewKeyAuth() PluginKeyAuthConfig {
	newKeyAuth := PluginKeyAuthConfig{}

	return newKeyAuth
}

func (keyAuthConfig PluginKeyAuthConfig) PluginConfigDefault() interface{} {
	pluginKeyAuth := PluginKeyAuth{
		Secret: utils.Md5(strconv.Itoa(int(time.Now().UnixNano()))),
	}

	return pluginKeyAuth
}

func (keyAuthConfig PluginKeyAuthConfig) PluginConfigParse(configInfo interface{}) (pluginKeyAuthConfig interface{}, err error) {

	pluginKeyAuth := PluginKeyAuth{
		Secret: "",
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

	err = json.Unmarshal(configInfoJson, &pluginKeyAuth)
	if err != nil {
		return
	}

	pluginKeyAuthConfig = pluginKeyAuth

	return
}

func (keyAuthConfig PluginKeyAuthConfig) PluginConfigCheck(configInfo interface{}) error {
	keyAuth, err := keyAuthConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginKeyAuth := keyAuth.(PluginKeyAuth)

	return keyAuthConfig.configValidator(pluginKeyAuth)
}

func (keyAuthConfig PluginKeyAuthConfig) configValidator(config PluginKeyAuth) error {

	fmt.Println("-----------", len(config.Secret), "-------", reflect.TypeOf(config.Secret), "-----------")

	if len(config.Secret) == 0 {
		return errors.New(fmt.Sprintf(
			keyAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.secret", "string"))
	}

	if len(config.Secret) < 10 {

		fmt.Println("==============")

		return errors.New(fmt.Sprintf(
			keyAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min"],
			"config.secret", 10))
	}

	if len(config.Secret) > 32 {
		return errors.New(fmt.Sprintf(
			keyAuthValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max"],
			"config.secret", 32))
	}

	return nil
}
