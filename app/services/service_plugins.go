package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type PluginsService struct {
}

var (
	pluginsService *PluginsService
	pluginsOnce    sync.Once
)

func NewPluginsService() *PluginsService {

	pluginsOnce.Do(func() {
		pluginsService = &PluginsService{}
	})

	return pluginsService
}

type PluginInfoService struct {
	ResID       string      `json:"res_id"`
	Name        string      `json:"name"`
	Key         string      `json:"key"`
	Icon        string      `json:"icon"`
	Type        int         `json:"type"`
	Description string      `json:"description"`
	Config      interface{} `json:"config"`
}

func (s *PluginsService) PluginInfoByResId(resId string) (PluginInfoService, error) {
	pluginInfo := PluginInfoService{}

	pluginModel := models.Plugins{}
	plugin, err := pluginModel.PluginInfoByResId(resId)

	if err != nil {
		return pluginInfo, err
	}
	if plugin.ID == 0 {
		return pluginInfo, errors.New(enums.CodeMessages(enums.PluginNull))
	}

	pluginContext, err := plugins.NewPluginContext(plugin.PluginKey)
	if err != nil {
		return pluginInfo, err
	}

	pluginConfig := pluginContext.StrategyPluginFormatDefault()
	pluginInfo = PluginInfoService{
		ResID:       plugin.ResID,
		Key:         plugin.PluginKey,
		Icon:        plugin.Icon,
		Type:        plugin.Type,
		Description: plugin.Description,
		Config:      pluginConfig,
	}

	return pluginInfo, nil
}

type PluginConfigListItem struct {
	ResID             string      `json:"res_id"`
	Name              string      `json:"name"`
	PluginKey         string      `json:"plugin_key"`
	PluginType        int         `json:"plugin_type"`
	PluginDescription string      `json:"plugin_description"`
	Enable            int         `json:"enable"`
	Config            interface{} `json:"config"`
}
type PluginConfigListResponse struct {
	List []PluginConfigListItem
}

func (s *PluginsService) PluginConfigList(resType int, resId string) (*PluginConfigListResponse, error) {

	res := &PluginConfigListResponse{
		List: []PluginConfigListItem{},
	}
	pluginConfigs, err := (&models.PluginConfigs{}).PluginConfigList(resType, resId)

	if err != nil {
		return res, err
	}

	if len(pluginConfigs) == 0 {
		return res, nil
	}

	pluginResID := []string{}
	for _, v := range pluginConfigs {
		pluginResID = append(pluginResID, v.PluginResID)
	}

	pluginList, err := (&models.Plugins{}).PluginInfosByResIds(pluginResID)

	if err != nil {
		return res, err
	}

	pluginMap := map[string]models.Plugins{}
	for k, v := range pluginList {
		pluginMap[v.ResID] = pluginList[k]
	}

	list := []PluginConfigListItem{}
	for _, v := range pluginConfigs {

		plugin := models.Plugins{}

		if tmp, ok := pluginMap[v.PluginResID]; ok {
			plugin = tmp
		}

		item := PluginConfigListItem{
			ResID:             v.ResID,
			Name:              v.Name,
			PluginKey:         plugin.PluginKey,
			PluginType:        plugin.Type,
			PluginDescription: plugin.Description,
			Enable:            v.Enable,
		}
		configContext, err := plugins.NewPluginContext(plugin.PluginKey)

		if err == nil {
			item.Config, _ = configContext.StrategyPluginParse(v.Config)
		}

		list = append(list, item)
	}

	return &PluginConfigListResponse{
		List: list,
	}, nil
}

func (s *PluginsService) PluginConfigInfoByResId(resId string) (*PluginConfigListItem, error) {

	pluginConfig, err := (&models.PluginConfigs{}).PluginConfigInfoByResId(resId)

	if err != nil {
		return &PluginConfigListItem{}, err
	}

	plugin, err := (&models.Plugins{}).PluginInfoByResId(pluginConfig.PluginResID)

	if err != nil {
		return &PluginConfigListItem{}, err
	}

	res := &PluginConfigListItem{
		ResID:             pluginConfig.ResID,
		Name:              pluginConfig.Name,
		PluginKey:         plugin.PluginKey,
		PluginType:        plugin.Type,
		PluginDescription: plugin.Description,
		Enable:            pluginConfig.Enable,
	}
	configContext, err := plugins.NewPluginContext(plugin.PluginKey)

	fmt.Println(err, plugin.PluginKey)
	if err == nil {
		res.Config, _ = configContext.StrategyPluginParse(pluginConfig.Config)
	}

	return res, nil
}

func (s *PluginsService) PluginConfigAdd(request *validators.ValidatorPluginConfigAdd) (string, error) {

	pluginInfo, err := s.PluginInfoByResId(request.PluginID)

	if err != nil {
		return "", errors.New(enums.CodeMessages(enums.PluginNull))
	}

	if request.Type == models.PluginConfigsTypeService {

		_, err := NewServicesService().ServiceInfoById(request.TargetID)
		if err != nil {
			return "", errors.New(enums.CodeMessages(enums.ServiceNull))
		}
	} else if request.Type == models.PluginConfigsTypeRouter {

		_, err := RouterInfoByResId(request.TargetID)
		if err != nil {
			return "", errors.New(enums.CodeMessages(enums.RouterNull))
		}
	}

	pluginContext, err := plugins.NewPluginContext(pluginInfo.Key)

	if err != nil {
		return "", err
	}

	if reflect.ValueOf(request.Config).IsNil() {
		request.Config = pluginContext.StrategyPluginFormatDefault()
	} else {
		err = pluginContext.StrategyPluginCheck(request.Config)

		if err != nil {
			return "", err
		}
		request.Config, _ = pluginContext.StrategyPluginParse(request.Config)
	}

	config, err := json.Marshal(request.Config)

	if err != nil {
		return "", err
	}

	pluginConfigResId, err := (&models.PluginConfigs{}).PluginConfigAdd(&models.PluginConfigs{
		Name:        request.Name,
		Type:        request.Type,
		TargetID:    request.TargetID,
		PluginResID: pluginInfo.ResID,
		PluginKey:   pluginInfo.Key,
		Config:      string(config),
		Enable:      request.Enable,
	})

	if err != nil {
		packages.Log.Error("create plugin config error")
		return pluginConfigResId, err
	}

	return pluginConfigResId, nil
}

func (s *PluginsService) PluginConfigUpdate(request *validators.ValidatorPluginConfigUpdate) error {

	pluginConfigInfo, err := (&models.PluginConfigs{}).PluginConfigInfoByResId(request.PluginConfigId)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.PluginConfigNull))
	}

	pluginInfo, err := s.PluginInfoByResId(pluginConfigInfo.PluginResID)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.PluginNull))
	}

	pluginContext, err := plugins.NewPluginContext(pluginInfo.Key)

	if err != nil {
		return err
	}

	if reflect.ValueOf(request.Config).IsNil() {
		request.Config = pluginContext.StrategyPluginFormatDefault()
	} else {
		err = pluginContext.StrategyPluginCheck(request.Config)

		if err != nil {
			return err
		}
		request.Config, _ = pluginContext.StrategyPluginParse(request.Config)
	}

	config, err := json.Marshal(request.Config)

	if err != nil {
		return err
	}

	err = (&models.PluginConfigs{}).PluginConfigUpdateColumns(
		pluginConfigInfo.ResID,
		pluginConfigInfo.Type,
		pluginConfigInfo.TargetID,
		map[string]interface{}{
			"name":   request.Name,
			"config": string(config),
		})

	if err != nil {
		packages.Log.Error("update plugin config error")
		return err
	}

	return nil
}

func (s *PluginsService) PluginConfigSwitchEnable(pluginConfigId string, enable int) error {

	pluginConfigInfo, err := (&models.PluginConfigs{}).PluginConfigInfoByResId(pluginConfigId)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.PluginConfigNull))
	}

	err = (&models.PluginConfigs{}).PluginConfigUpdateColumns(
		pluginConfigInfo.ResID,
		pluginConfigInfo.Type,
		pluginConfigInfo.TargetID,
		map[string]interface{}{
			"enable": enable,
		})

	if err != nil {
		packages.Log.Error("Failed to update plugin enable switch")
		return err
	}

	return nil
}

func (s *PluginsService) PluginConfigDelete(pluginConfigId string) error {

	pluginConfigInfo, err := (&models.PluginConfigs{}).PluginConfigInfoByResId(pluginConfigId)

	if err != nil {
		return errors.New(enums.CodeMessages(enums.PluginConfigNull))
	}

	err = (&models.PluginConfigs{}).PluginConfigDelete(
		pluginConfigInfo.ResID,
		pluginConfigInfo.Type,
		pluginConfigInfo.TargetID,
	)

	if err != nil {
		packages.Log.Error("Failed to delete plugin config")
		return err
	}

	return nil
}
