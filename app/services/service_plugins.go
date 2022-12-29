package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/rpc"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
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
	if plugin.ResID == "" {
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
	Icon              string      `json:"icon"`
	PluginKey         string      `json:"plugin_key"`
	PluginType        int         `json:"plugin_type"`
	PluginDescription string      `json:"plugin_description"`
	Enable            int         `json:"enable"`
	Config            interface{} `json:"config"`
}
type PluginConfigListResponse struct {
	List []PluginConfigListItem `json:"list"`
}

func (s *PluginsService) PluginConfigList(resType int, resId string) (*PluginConfigListResponse, error) {

	res := &PluginConfigListResponse{
		List: []PluginConfigListItem{},
	}
	pluginConfigs, err := (&models.PluginConfigs{}).PluginConfigList(packages.GetDb(), resType, resId, 0)

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
			Icon:              plugin.Icon,
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

	if err == nil {
		res.Config, _ = configContext.StrategyPluginParse(pluginConfig.Config)
	}

	return res, nil
}

func (s *PluginsService) PluginConfigAdd(request *validators.ValidatorPluginConfigAdd) (pluginConfigResId string, err error) {

	var pluginInfo PluginInfoService
	pluginInfo, err = s.PluginInfoByResId(request.PluginID)

	if err != nil {
		err = errors.New(enums.CodeMessages(enums.PluginNull))
		return
	}

	if request.Type == models.PluginConfigsTypeService {

		_, err = NewServicesService().ServiceInfoById(request.TargetID)
		if err != nil {
			err = errors.New(enums.CodeMessages(enums.ServiceNull))
			return
		}
	} else if request.Type == models.PluginConfigsTypeRouter {

		_, err = RouterInfoByResId(request.TargetID)
		if err != nil {
			err = errors.New(enums.CodeMessages(enums.RouterNull))
			return
		}
	}

	var pluginContext plugins.PluginContext
	pluginContext, err = plugins.NewPluginContext(pluginInfo.Key)

	if err != nil {
		return
	}

	err = pluginContext.StrategyPluginCheck(request.Config)
	if err != nil {
		return
	}

	request.Config, _ = pluginContext.StrategyPluginParse(request.Config)
	config, err := json.Marshal(request.Config)

	if err != nil {
		return
	}

	pluginConfigResId, err = (&models.PluginConfigs{}).PluginConfigAdd(&models.PluginConfigs{
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
		return
	}

	return
}

func (s *PluginsService) PluginConfigUpdate(request *validators.ValidatorPluginConfigUpdate) error {

	pluginConfigInfo, err := (&models.PluginConfigs{}).PluginConfigInfoByResId(request.PluginConfigId)

	if err != nil {
		return err
	}

	if pluginConfigInfo.ResID == "" {
		return errors.New(enums.CodeMessages(enums.PluginConfigNull))
	}

	pluginInfo, err := s.PluginInfoByResId(pluginConfigInfo.PluginResID)

	if err != nil {
		return err
	}

	if pluginInfo.ResID == "" {
		return errors.New(enums.CodeMessages(enums.PluginNull))
	}

	pluginContext, err := plugins.NewPluginContext(pluginInfo.Key)

	if err != nil {
		return err
	}

	request.Config, _ = pluginContext.StrategyPluginParse(request.Config)
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

func SyncPluginToDataSide(tx *gorm.DB, resType int, targetId string) ([]models.PluginConfigs, error) {
	// 获取控制面已绑定插件，同步至远程
	pluginConfigList, err := (&models.PluginConfigs{}).PluginConfigList(tx, resType, targetId, utils.EnableOn)

	if err != nil {
		return []models.PluginConfigs{}, err
	}

	if len(pluginConfigList) == 0 {
		return []models.PluginConfigs{}, nil
	}
	success := []models.PluginConfigs{}
	// 同步服务插件
	for k, v := range pluginConfigList {

		pluginContext, err := plugins.NewPluginContext(v.PluginKey)

		if err != nil {
			continue
		}

		config, err := pluginContext.StrategyPluginParse(v.Config)

		if err != nil {
			continue
		}
		pluginPutRequest := &rpc.PluginPutRequest{
			Name:   v.ResID,
			Key:    v.PluginKey,
			Config: config,
		}
		err = rpc.NewApiOak().PluginPut(pluginPutRequest)

		if err != nil {
			continue
		}

		success = append(success, pluginConfigList[k])
	}

	return success, nil
}

func PluginBasicInfoMaintain() {

	pluginModel := models.Plugins{}
	dbPluginList, _ := pluginModel.PluginAllList()

	dbPluginMapResId := make(map[string]models.Plugins)

	for _, dbPluginInfo := range dbPluginList {
		dbPluginMapResId[dbPluginInfo.ResID] = dbPluginInfo
	}

	configPluginList := utils.AllConfigPluginData()

	for _, configPluginInfo := range configPluginList {

		dbPluginMapInfo, ok := dbPluginMapResId[configPluginInfo.ResID]

		if ok {
			if (configPluginInfo.PluginKey != dbPluginMapInfo.PluginKey) ||
				(configPluginInfo.Type != dbPluginMapInfo.Type) {

				dbPluginMapInfo.PluginKey = configPluginInfo.PluginKey
				dbPluginMapInfo.Type = configPluginInfo.Type
				dbPluginMapInfo.Icon = configPluginInfo.Icon
				dbPluginMapInfo.Description = configPluginInfo.Description
				_ = pluginModel.PluginUpdate(configPluginInfo.ResID, &dbPluginMapInfo)
			}

		} else {

			_ = pluginModel.PluginDelByPluginKeys([]string{configPluginInfo.PluginKey}, []string{})

			newPluginData := pluginModel
			newPluginData.ResID = configPluginInfo.ResID
			newPluginData.Type = configPluginInfo.Type
			newPluginData.PluginKey = configPluginInfo.PluginKey
			newPluginData.Icon = configPluginInfo.Icon
			newPluginData.Description = configPluginInfo.Description

			_ = pluginModel.PluginAdd(&newPluginData)
		}
	}
}

type PluginConfigDefault struct {
	ResID       string      `json:"res_id"`
	PluginKey   string      `json:"plugin_key"`
	Icon        string      `json:"icon"`
	Type        int         `json:"type"`
	Description string      `json:"description"`
	Config      interface{} `json:"config"`
}

func (s *PluginsService) PluginConfigDefault(pluginResId string) (pluginConfigDefault PluginConfigDefault, err error) {
	pluginModel := models.Plugins{}
	var pluginInfo models.Plugins
	pluginInfo, err = pluginModel.PluginInfoByResId(pluginResId)
	if err != nil {
		return
	}

	if pluginInfo.ResID == "" {
		return
	}

	var pluginContext plugins.PluginContext
	pluginContext, err = plugins.NewPluginContext(pluginInfo.PluginKey)

	if err != nil {
		return
	}

	pluginConfigDefault.ResID = pluginInfo.ResID
	pluginConfigDefault.PluginKey = pluginInfo.PluginKey
	pluginConfigDefault.Icon = pluginInfo.Icon
	pluginConfigDefault.Type = pluginInfo.Type
	pluginConfigDefault.Description = pluginInfo.Description
	pluginConfigDefault.Config = pluginContext.StrategyPluginFormatDefault()

	return
}
