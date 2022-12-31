package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var corsValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required":   "[%s] is a required field,expected type: %s",
		"max_length": "[%s] length must be less than or equal to %d",
		"min_length": "[%s] length must be greater than or equal to %d",
		"max_number": "[%s] must be %d or less",
		"min_number": "[%s] must be %d or greater",
		"oneOf":      "[%s] must be a value that exists in [%s]",
	},
	utils.LocalZh: {
		"required":   "[%s]为必填字段，期望类型:%s",
		"max_length": "[%s]长度必须小于或等于%d",
		"min_length": "[%s]长度必须大于或等于%d",
		"max_number": "[%s]必须小于或等于%d",
		"min_number": "[%s]必须大于或等于%d",
		"oneOf":      "[%s]必须是存在于[%s]中的值",
	},
}

var allMethodsList = []string{
	"*",
	"GET",
	"PUT",
	"POST",
	"HEAD",
	"PATCH",
	"TRACE",
	"DELETE",
	"OPTIONS",
	"CONNECT",
}

type PluginCorsConfig struct{}

type PluginCors struct {
	AllowMethods    string `json:"allow_methods"`
	AllowOrigins    string `json:"allow_origins"`
	AllowHeaders    string `json:"allow_headers"`
	MaxAge          int    `json:"max_age"`
	AllowCredential bool   `json:"allow_credential"`
}

func NewCors() PluginCorsConfig {
	newCors := PluginCorsConfig{}

	return newCors
}

func (corsConfig PluginCorsConfig) PluginConfigDefault() interface{} {
	pluginCors := PluginCors{
		AllowMethods:    "*",
		AllowOrigins:    "*",
		AllowHeaders:    "*",
		MaxAge:          0,
		AllowCredential: false,
	}

	return pluginCors
}

func (corsConfig PluginCorsConfig) PluginConfigParse(configInfo interface{}) (pluginCorsConfig interface{}, err error) {

	pluginCors := PluginCors{
		AllowMethods:    "",
		AllowOrigins:    "",
		AllowHeaders:    "",
		MaxAge:          0,
		AllowCredential: false,
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

	err = json.Unmarshal(configInfoJson, &pluginCors)
	if err != nil {
		return
	}

	pluginCors.AllowMethods = strings.TrimSpace(pluginCors.AllowMethods)
	pluginCors.AllowOrigins = strings.TrimSpace(pluginCors.AllowOrigins)
	pluginCors.AllowHeaders = strings.TrimSpace(pluginCors.AllowHeaders)

	if len(pluginCors.AllowMethods) > 0 {
		allowMethodsArr := strings.Split(pluginCors.AllowMethods, ",")
		for key, allowMethodsArrInfo := range allowMethodsArr {
			allowMethodsArr[key] = strings.ToUpper(strings.TrimSpace(allowMethodsArrInfo))
		}
		pluginCors.AllowMethods = strings.Join(allowMethodsArr, ",")
	}

	pluginCorsConfig = pluginCors

	return
}

func (corsConfig PluginCorsConfig) PluginConfigCheck(configInfo interface{}) error {
	cors, err := corsConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginCors := cors.(PluginCors)

	return corsConfig.configValidator(pluginCors)
}

func (corsConfig PluginCorsConfig) configValidator(config PluginCors) error {

	if len(config.AllowOrigins) > 80 {
		return errors.New(fmt.Sprintf(
			corsValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_length"],
			"config.allow_origins", 80))
	}

	if len(config.AllowHeaders) > 80 {
		return errors.New(fmt.Sprintf(
			corsValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_length"],
			"config.allow_headers", 80))
	}

	if config.MaxAge < 0 {
		return errors.New(fmt.Sprintf(
			corsValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.max_age", 0))
	}

	if config.MaxAge > 86400 {
		return errors.New(fmt.Sprintf(
			corsValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.max_age", 86400))
	}

	if len(config.AllowMethods) > 0 {

		allowMethodsArr := strings.Split(config.AllowMethods, ",")

		if len(allowMethodsArr) > 0 {

			allMethodsListMap := make(map[string]byte)
			for _, allMethodsInfo := range allMethodsList {
				_, ok := allMethodsListMap[allMethodsInfo]
				if !ok {
					allMethodsListMap[allMethodsInfo] = 0
				}
			}

			for _, allowMethodsArrInfo := range allowMethodsArr {

				_, ok := allMethodsListMap[allowMethodsArrInfo]

				if !ok {
					return errors.New(fmt.Sprintf(
						corsValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
						"config.allow_methods", strings.Join(allMethodsList, ",")))
				}
			}
		}
	}

	return nil
}
