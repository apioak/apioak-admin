package plugins

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var mockValidatorErrorMessages = map[string]map[string]string{
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

var responseTypeList = []string{
	"application/json",
	"text/html",
	"text/xml",
}

type PluginMockConfig struct{}

type PluginMock struct {
	ResponseType string            `json:"response_type"`
	HttpCode     int               `json:"http_code"`
	HttpBody     string            `json:"http_body"`
	HttpHeaders  map[string]string `json:"http_headers"`
}

func NewMock() PluginMockConfig {
	newMock := PluginMockConfig{}

	return newMock
}

func (mockConfig PluginMockConfig) PluginConfigDefault() interface{} {
	pluginMock := PluginMock{
		ResponseType: "application/json",
		HttpCode:     0,
		HttpBody:     "",
		HttpHeaders: map[string]string{},
	}

	return pluginMock
}

func (mockConfig PluginMockConfig) PluginConfigParse(configInfo interface{}) (pluginMockConfig interface{}, err error) {

	pluginMock := PluginMock{
		ResponseType: "application/json",
		HttpCode:     -999,
		HttpBody:     "",
		HttpHeaders: map[string]string{},
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

	err = json.Unmarshal(configInfoJson, &pluginMock)
	if err != nil {
		return
	}

	pluginMock.ResponseType = strings.TrimSpace(pluginMock.ResponseType)
	pluginMock.HttpBody = strings.TrimSpace(pluginMock.HttpBody)

	pluginMockConfig = pluginMock

	return
}

func (mockConfig PluginMockConfig) PluginConfigCheck(configInfo interface{}) error {
	mock, err := mockConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginMock := mock.(PluginMock)

	return mockConfig.configValidator(pluginMock)
}

func (mockConfig PluginMockConfig) configValidator(config PluginMock) error {

	if config.HttpCode == -999 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.http_code", "int"))
	}

	if len(config.HttpBody) == 0 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.http_body", "string"))
	}

	if len(config.ResponseType) == 0 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.response_type", "string"))
	}

	if config.HttpCode < 100 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.http_code", 100))
	}

	if config.HttpCode > 599 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.http_code", 599))
	}

	if len(config.HttpBody) < 1 {
		return errors.New(fmt.Sprintf(
			mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_length"],
			"config.http_body", 1))
	}

	if len(config.ResponseType) > 0 {

		responseTypeListMap := make(map[string]byte)
		for _, responseTypeListInfo := range responseTypeList {
			_, ok := responseTypeListMap[responseTypeListInfo]
			if !ok {
				responseTypeListMap[responseTypeListInfo] = 0
			}
		}

		_, exist := responseTypeListMap[config.ResponseType]

		if !exist {
			return errors.New(fmt.Sprintf(
				mockValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
				"config.response_type", strings.Join(responseTypeList, " ")))
		}
	}

	return nil
}
