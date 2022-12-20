package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/rpc"
	"apioak-admin/app/services/plugins"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func CheckRouterExist(routerResId string, serviceResId string) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerResId, serviceResId)

	if len(routerInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterNull))
	}

	return nil
}

func CheckRouterRelease(routerResId string) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerResId, "")

	if len(routerInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterNull))
	}

	if routerInfo.Release == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func CheckServiceRouterPath(path string) error {
	if path == utils.DefaultRouterPath {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathNoPermission))
	}

	if strings.Index(path, utils.DefaultRouterPath) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathForbiddenPrefix))
	}

	return nil
}

func CheckExistServiceRouterPath(serviceResId string, path string, filterRouterResIds []string) error {
	routerModel := models.Routers{}
	routerPaths, err := routerModel.RouterInfosByServiceRouterPath(serviceResId, []string{path}, filterRouterResIds)
	if err != nil {
		return err
	}

	if len(routerPaths) == 0 {
		return nil
	}

	existRouterPath := make([]string, 0)
	tmpExistRouterPathMap := make(map[string]byte, 0)
	for _, routerPath := range routerPaths {
		_, exist := tmpExistRouterPathMap[routerPath.RouterPath]
		if exist {
			continue
		}

		existRouterPath = append(existRouterPath, routerPath.RouterPath)
		tmpExistRouterPathMap[routerPath.RouterPath] = 0
	}

	if len(existRouterPath) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.RouterPathExist), strings.Join(existRouterPath, ",")))
	}

	return nil
}

func CheckRouterEnableChange(routerId string, enable int) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerId, "")

	if routerInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func RouterCreate(routerData *validators.ValidatorRouterAddUpdate) (routerResId string, err error) {
	createRouterData := models.Routers{
		ServiceResID:   routerData.ServiceResID,
		UpstreamResID:  routerData.UpstreamResID,
		RouterName:     routerData.RouterName,
		RequestMethods: routerData.RequestMethods,
		RouterPath:     routerData.RouterPath,
		Enable:         routerData.Enable,
		Release:        utils.ReleaseStatusU,
	}

	createUpstreamData := models.Upstreams{
		Algorithm:      routerData.LoadBalance,
		ConnectTimeout: routerData.ConnectTimeout,
		WriteTimeout:   routerData.WriteTimeout,
		ReadTimeout:    routerData.ReadTimeout,
	}

	createUpstreamNodes := make([]models.UpstreamNodes, 0)
	if len(routerData.UpstreamNodes) > 0 {
		for _, upstreamNode := range routerData.UpstreamNodes {
			var ipType string
			ipType, err = utils.DiscernIP(upstreamNode.NodeIp)
			if err != nil {
				return
			}
			ipNameIdMap := utils.IpNameIdMap()

			createUpstreamNodes = append(createUpstreamNodes, models.UpstreamNodes{
				NodeIP:      upstreamNode.NodeIp,
				IPType:      ipNameIdMap[ipType],
				NodePort:    upstreamNode.NodePort,
				NodeWeight:  upstreamNode.NodeWeight,
				Health:      upstreamNode.Health,
				HealthCheck: upstreamNode.HealthCheck,
			})
		}
	}

	routerResId, err = createRouterData.RouterAdd(createRouterData, createUpstreamData, createUpstreamNodes)

	if err != nil {
		return
	}

	return
}

type routerPlugin struct {
	ResID  string `json:"res_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Icon   string `json:"icon"`
	Type   int    `json:"type"`
	Enable int    `json:"enable"`
}

type RouterListItem struct {
	ResID          string         `json:"res_id"`
	ServiceResID   string         `json:"service_res_id"`
	ServiceName    string         `json:"service_name"`
	RouterName     string         `json:"router_name"`
	RequestMethods []string       `json:"request_methods"`
	RouterPath     string         `json:"router_path"`
	Enable         int            `json:"enable"`
	Release        int            `json:"release"`
	PluginList     []routerPlugin `json:"plugin_list"`
}

func (s *RouterListItem) RouterListPage(serviceResId string, param *validators.ValidatorRouterList) (
	routerList []RouterListItem, total int, err error) {

	routerModel := models.Routers{}
	routerInfos := make([]models.Routers, 0)
	routerInfos, total, err = routerModel.RouterListPage(serviceResId, param)

	routerServiceResIds := make([]string, 0)
	routerServiceResIdsMap := make(map[string]byte)

	routerResIds := make([]string, 0)
	if len(routerInfos) != 0 {
		for _, routerInfo := range routerInfos {
			routerListItem := RouterListItem{}
			routerListItem.ResID = routerInfo.ResID
			routerListItem.ServiceResID = routerInfo.ServiceResID
			routerListItem.RouterName = routerInfo.RouterName
			routerListItem.RequestMethods = strings.Split(routerInfo.RequestMethods, ",")
			routerListItem.RouterPath = routerInfo.RouterPath
			routerListItem.Enable = routerInfo.Enable
			routerListItem.Release = routerInfo.Release
			routerListItem.PluginList = make([]routerPlugin, 0)
			routerList = append(routerList, routerListItem)
			routerResIds = append(routerResIds, routerInfo.ResID)

			if _, ok := routerServiceResIdsMap[routerInfo.ServiceResID]; ok == false {
				routerServiceResIds = append(routerServiceResIds, routerInfo.ServiceResID)
			}
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, routerResIds)
	if err != nil {
		return
	}

	if len(pluginConfigList) > 0 {

		pluginResIds := make([]string, 0)
		pluginResIdsMap := make(map[string]byte)
		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginResIdsMap[pluginConfigInfo.PluginResID]
			if ok == false {
				pluginResIds = append(pluginResIds, pluginConfigInfo.PluginResID)
			}
		}

		pluginModel := models.Plugins{}
		pluginList := make([]models.Plugins, 0)
		pluginList, err = pluginModel.PluginAllList()
		if err != nil {
			return
		}

		pluginListMap := make(map[string]models.Plugins)
		for _, pluginInfo := range pluginList {
			pluginListMap[pluginInfo.ResID] = pluginInfo
		}

		pluginConfigMapList := make(map[string][]routerPlugin)
		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginConfigMapList[pluginConfigInfo.TargetID]
			if ok == false {
				pluginConfigMapList[pluginConfigInfo.TargetID] = make([]routerPlugin, 0)
			}
			pluginConfigMapList[pluginConfigInfo.TargetID] = append(pluginConfigMapList[pluginConfigInfo.TargetID], routerPlugin{
				ResID:  pluginConfigInfo.ResID,
				Name:   pluginConfigInfo.Name,
				Key:    pluginConfigInfo.PluginKey,
				Enable: pluginConfigInfo.Enable,
				Icon:   pluginListMap[pluginConfigInfo.PluginResID].Icon,
				Type:   pluginListMap[pluginConfigInfo.PluginResID].Type,
			})
		}

		if len(routerList) > 0 {
			for key, routerInfo := range routerList {
				routerPluginList, ok := pluginConfigMapList[routerInfo.ResID]
				if ok {
					routerList[key].PluginList = routerPluginList
				}
			}
		}
	}

	if len(routerServiceResIds) > 0 {

		serviceModel := models.Services{}
		serviceList := make([]models.Services, 0)
		serviceList, err = serviceModel.ServiceListByResIds(routerServiceResIds)
		if err != nil {
			return
		}

		serviceMap := make(map[string]models.Services)
		for _, serviceInfo := range serviceList {
			serviceMap[serviceInfo.ResID] = serviceInfo
		}

		if len(routerList) > 0 {
			for key, routerInfo := range routerList {
				serviceInfo, ok := serviceMap[routerInfo.ServiceResID]
				if ok {
					routerList[key].ServiceName = serviceInfo.Name
				}
			}
		}
	}

	return
}

type StructRouterInfo struct {
	ID             string             `json:"id"`
	ServiceID      string             `json:"service_id"`
	RouterName     string             `json:"router_name"`
	RequestMethods []string           `json:"request_methods"`
	RouterPath     string             `json:"router_path"`
	Enable         int                `json:"enable"`
	Release        int                `json:"release"`
	Upstream       UpstreamItem       `json:"upstream"`
	UpstreamNodes  []UpstreamNodeItem `json:"upstream_nodes"`
}

func (s *StructRouterInfo) RouterInfoByServiceRouterId(serviceResId string, routerResId string) (routerDetail StructRouterInfo, err error) {
	routerModel := &models.Routers{}
	routerModelDetail, routerModelDetailErr := routerModel.RouterInfosByServiceRouterId(serviceResId, routerResId)
	if routerModelDetailErr != nil {
		err = routerModelDetailErr
		return
	}

	routerDetail.ID = routerModelDetail.ResID
	routerDetail.ServiceID = routerModelDetail.ServiceResID
	routerDetail.RouterName = routerModelDetail.RouterName
	routerDetail.RequestMethods = strings.Split(routerModelDetail.RequestMethods, ",")
	routerDetail.RouterPath = routerModelDetail.RouterPath
	routerDetail.Enable = routerModelDetail.Enable
	routerDetail.Release = routerModelDetail.Release

	upstreamItem := UpstreamItem{}
	upstreamDetail, upstreamDetailErr := upstreamItem.UpstreamDetailByResId(routerModelDetail.UpstreamResID)
	if upstreamDetailErr == nil {
		routerDetail.Upstream = upstreamDetail
	}

	upstreamNodeItem := UpstreamNodeItem{}
	upstreamNodeList, upstreamNodeListErr := upstreamNodeItem.UpstreamNodeListByUpstreamResIds([]string{routerModelDetail.UpstreamResID})
	if upstreamNodeListErr == nil {
		routerDetail.UpstreamNodes = upstreamNodeList
	}

	return
}

func RouterUpdate(routerResId string, routerData validators.ValidatorRouterAddUpdate) (err error) {
	routerModel := models.Routers{}

	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	updateRouterData := make(map[string]interface{})
	updateRouterData["request_methods"] = routerData.RequestMethods
	updateRouterData["router_path"] = routerData.RouterPath
	updateRouterData["enable"] = routerData.Enable

	if len(routerData.RouterName) != 0 {
		updateRouterData["router_name"] = routerData.RouterName
	}
	if routerDetail.Release == utils.ReleaseStatusY {
		updateRouterData["release"] = utils.ReleaseStatusT
	}

	var upstreamResId string
	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {

		upstreamModel := models.Upstreams{}
		upstreamNodeModel := models.UpstreamNodes{}

		if len(routerData.UpstreamNodes) == 0 {
			err = UpstreamRelease([]string{upstreamResId}, utils.ReleaseTypeDelete)
			if err != nil {
				return
			}

			if err = tx.Table(upstreamModel.TableName()).
				Where("res_id = ?", routerDetail.UpstreamResID).
				Delete(&upstreamModel).Error; err != nil {
				return
			}

			if err = tx.Table(upstreamNodeModel.TableName()).
				Where("upstream_res_id = ?", routerDetail.UpstreamResID).
				Delete(&upstreamNodeModel).Error; err != nil {
				return
			}

			upstreamResId = routerDetail.UpstreamResID
			updateRouterData["upstream_res_id"] = ""
		} else {
			var upstreamInfo models.Upstreams
			upstreamInfo, err = upstreamInfo.UpstreamDetailByResId(routerDetail.UpstreamResID)
			if err != nil {
				return
			}

			if len(upstreamInfo.ResID) == 0 {
				upstreamResId, err = upstreamModel.ModelUniqueId()
				if err != nil {
					return
				}

				if err = tx.Table(upstreamModel.TableName()).
					Create(&models.Upstreams{
						ResID:          upstreamResId,
						Name:           upstreamResId,
						Algorithm:      routerData.LoadBalance,
						ReadTimeout:    routerData.ReadTimeout,
						WriteTimeout:   routerData.WriteTimeout,
						ConnectTimeout: routerData.ConnectTimeout,
					}).Error; err != nil {
					return
				}

				updateRouterData["upstream_res_id"] = upstreamResId
			} else {
				if err = tx.Table(upstreamModel.TableName()).
					Where("res_id = ?", upstreamInfo.ResID).
					Updates(models.Upstreams{
						Algorithm:      routerData.LoadBalance,
						ReadTimeout:    routerData.ReadTimeout,
						WriteTimeout:   routerData.WriteTimeout,
						ConnectTimeout: routerData.ConnectTimeout,
					}).Error; err != nil {
					return
				}

				upstreamResId = routerDetail.UpstreamResID
				updateRouterData["upstream_res_id"] = routerDetail.UpstreamResID
			}

			addNodeList, updateNodeList, delNodeResIds := DiffUpstreamNode(upstreamResId, routerData.UpstreamNodes)

			if len(addNodeList) > 0 {
				if err = tx.Create(&addNodeList).Error; err != nil {
					return
				}
			}

			if len(updateNodeList) > 0 {
				for _, updateNodeInfo := range updateNodeList {
					if err = tx.Table(upstreamNodeModel.TableName()).
						Where("res_id = ?", updateNodeInfo.ResID).
						Updates(&updateNodeInfo).Error; err != nil {
						return
					}
				}
			}

			if len(delNodeResIds) > 0 {
				if err = tx.Table(upstreamNodeModel.TableName()).
					Where("res_id in ?", delNodeResIds).
					Delete(&upstreamNodeModel).Error; err != nil {
					return
				}
			}
		}

		if err = tx.Table(routerModel.TableName()).
			Where("res_id = ?", routerResId).
			Updates(&updateRouterData).Error; err != nil {
			return
		}

		return
	})

	if err != nil {
		return
	}

	err = UpstreamRelease([]string{upstreamResId}, utils.ReleaseTypePush)
	if err != nil {
		return
	}

	return
}

func filterPushedServiceRouterResIds(routerResIds []string) (opRoutersResIds []string, publishedRouterResIds []string, err error) {
	if len(routerResIds) == 0 {
		return
	}

	routerModel := models.Routers{}
	routerList := make([]models.Routers, 0)
	routerList, err = routerModel.RouterListByRouterResIds(routerResIds)
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	serviceResIds := make([]string, 0)

	for _, routerInfo := range routerList {
		if len(routerInfo.ServiceResID) > 0 {
			serviceResIds = append(serviceResIds, routerInfo.ServiceResID)

			if routerInfo.Release != utils.ReleaseStatusU {
				publishedRouterResIds = append(publishedRouterResIds, routerInfo.ResID)
			}
		}
	}

	serviceModel := models.Services{}
	serviceList := make([]models.Services, 0)
	serviceList, err = serviceModel.ServiceListByResIds(serviceResIds)
	if err != nil {
		return
	}

	publishedServiceResIdsMap := make(map[string]byte)
	for _, serviceInfo := range serviceList {
		if serviceInfo.Release != utils.ReleaseStatusU {
			publishedServiceResIdsMap[serviceInfo.ResID] = 0
		}
	}

	for _, routerInfo := range routerList {
		_, ok := publishedServiceResIdsMap[routerInfo.ServiceResID]
		if ok {
			opRoutersResIds = append(opRoutersResIds, routerInfo.ResID)
		}
	}

	return
}

func RouterUpstreamRelease(routerResIds []string, releaseType string) (err error) {
	if len(routerResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	routerModel := models.Routers{}
	var routerList []models.Routers
	routerList, err = routerModel.RouterListByRouterResIds(routerResIds)
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	serviceResIds := make([]string, 0)

	for _, routerInfo := range routerList {
		if len(routerInfo.ServiceResID) > 0 {
			serviceResIds = append(serviceResIds, routerInfo.ServiceResID)
		}
	}

	serviceModel := models.Services{}
	serviceList := make([]models.Services, 0)
	serviceList, err = serviceModel.ServiceListByResIds(serviceResIds)
	if err != nil {
		return
	}

	publishedServiceResIdsMap := make(map[string]byte)
	for _, serviceInfo := range serviceList {
		if serviceInfo.Release != utils.ReleaseStatusU {
			publishedServiceResIdsMap[serviceInfo.ResID] = 0
		}
	}

	toBeOpUpstreamResIds := make([]string, 0)
	toBeOpRouterList := make([]models.Routers, 0)
	for _, routerInfo := range routerList {

		// _, ok := publishedServiceResIdsMap[routerInfo.ServiceResID]
		// if !ok {
		// 	continue
		// }

		toBeOpRouterList = append(toBeOpRouterList, routerInfo)

		if len(routerInfo.UpstreamResID) > 0 {
			toBeOpUpstreamResIds = append(toBeOpUpstreamResIds, routerInfo.UpstreamResID)
		}
	}

	if len(toBeOpRouterList) == 0 {
		return
	}

	routerConfigList := make([]rpc.RouterConfig, 0)
	routerConfigResIds := make([]string, 0)
	for _, toBeOpRouterInfo := range toBeOpRouterList {
		var routerConfig rpc.RouterConfig
		routerConfig, err = generateRouterConfig(toBeOpRouterInfo)
		if err != nil {
			return
		}

		if len(routerConfig.Name) == 0 {
			continue
		}

		routerConfigList = append(routerConfigList, routerConfig)
		routerConfigResIds = append(routerConfigResIds, toBeOpRouterInfo.ResID)
	}

	pluginConfigModel := models.PluginConfigs{}
	routerPluginConfigList := make([]models.PluginConfigs, 0)
	routerPluginConfigList, err = pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, routerConfigResIds)
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf("---------%+v--------", routerPluginConfigList))

	newApiOak := rpc.NewApiOak()
	if releaseType == utils.ReleaseTypePush {

		if len(routerPluginConfigList) > 0 {
			for _, routerPluginConfigInfo := range routerPluginConfigList {

				var pluginContext plugins.PluginContext
				pluginContext, err = plugins.NewPluginContext(routerPluginConfigInfo.PluginKey)

				if err != nil {
					continue
				}

				var config interface{}
				config, err = pluginContext.StrategyPluginParse(routerPluginConfigInfo.Config)

				err = newApiOak.PluginPut(&rpc.PluginPutRequest{
					Name:   routerPluginConfigInfo.ResID,
					Key:    routerPluginConfigInfo.PluginKey,
					Config: config,
				})

				if err != nil {
					continue
				}
			}
		}

		err = newApiOak.RouterDelete(routerResIds)
		if err != nil {
			return
		}

		err = UpstreamRelease(toBeOpUpstreamResIds, releaseType)
		if err != nil {
			return
		}

		err = newApiOak.RouterPut(routerConfigList)
		if err != nil {
			return
		}

		for _, toBeOpRouterInfo := range toBeOpRouterList {
			err = routerModel.RouterSwitchRelease(toBeOpRouterInfo.ResID, utils.ReleaseStatusY)
			if err != nil {
				return
			}
		}

	} else {
		err = newApiOak.RouterDelete(routerResIds)
		if err != nil {
			return
		}

		err = UpstreamRelease(toBeOpUpstreamResIds, releaseType)
		if err != nil {
			return
		}

		apioakRouterConfigList := make([]rpc.RouterConfig, 0)
		apioakRouterConfigList, err = newApiOak.RouterGet(routerConfigResIds)
		if err != nil {
			return
		}

		if len(apioakRouterConfigList) > 0 {
			for _, apioakRouterConfigInfo := range apioakRouterConfigList {
				if len(apioakRouterConfigInfo.Plugins) == 0 {
					continue
				}

				for _, pluginIdName := range apioakRouterConfigInfo.Plugins {
					err = newApiOak.PluginDelete(pluginIdName.Id)
					if err != nil {
						return
					}
				}
			}
		}
	}

	return
}

func RouterRelease(routerResIds []string, releaseType string) (err error) {
	if len(routerResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	opRouterResIds := make([]string, 0)
	publishedRouterResIds := make([]string, 0)
	opRouterResIds, publishedRouterResIds, err = filterPushedServiceRouterResIds(routerResIds)

	routerModel := models.Routers{}
	routerList := make([]models.Routers, 0)
	routerList, err = routerModel.RouterListByRouterResIds(opRouterResIds)
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	newApiOak := rpc.NewApiOak()
	if releaseType == utils.ReleaseTypePush {

		routerConfigList := make([]rpc.RouterConfig, 0)
		for _, routerInfo := range routerList {
			var routerConfig rpc.RouterConfig
			routerConfig, err = generateRouterConfig(routerInfo)
			if err != nil {
				return
			}

			if len(routerConfig.Name) == 0 {
				continue
			}

			routerConfigList = append(routerConfigList, routerConfig)
		}

		err = newApiOak.RouterPut(routerConfigList)
		if err != nil {
			return
		}

	} else {
		err = newApiOak.RouterDelete(publishedRouterResIds)
	}

	return
}

func generateRouterConfig(routerInfo models.Routers) (rpc.RouterConfig, error) {
	routerConfig := rpc.RouterConfig{}

	routerConfig.Name = routerInfo.ResID
	routerConfig.Methods = strings.Split(routerInfo.RequestMethods, ",")
	routerConfig.Paths = append(routerConfig.Paths, routerInfo.RouterPath)
	routerConfig.Enabled = false
	if routerInfo.Enable == utils.EnableOn {
		routerConfig.Enabled = true
	}
	routerConfig.Headers = make(map[string]string)
	routerConfig.Service.Name = routerInfo.ServiceResID
	routerConfig.Upstream.Name = routerInfo.UpstreamResID
	routerConfig.Plugins = make([]rpc.ConfigObjectName, 0)

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, []string{routerInfo.ResID})
	if err != nil {
		return routerConfig, err
	}

	if len(pluginConfigList) > 0 {
		for _, pluginConfigInfo := range pluginConfigList {
			if pluginConfigInfo.Enable == utils.EnableOff {
				continue
			}

			routerConfig.Plugins = append(routerConfig.Plugins, rpc.ConfigObjectName{
				Name: pluginConfigInfo.ResID,
			})
		}
	}

	return routerConfig, nil
}

func CheckEditDefaultPathRouter(routerId string) error {
	routerModel := models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerId, "")
	if routerInfo.RouterPath == utils.DefaultRouterPath {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathNoPermission))
	}

	return nil
}

func RouterDelete(routerResId string) (err error) {
	routerModel := models.Routers{}

	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	if routerDetail.ResID != routerResId {
		return
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(routerModel.TableName()).
			Where("res_id = ?", routerResId).
			Delete(&routerModel).Error; err != nil {
			return
		}

		upstreamModel := models.Upstreams{}
		if err = tx.Table(upstreamModel.TableName()).
			Where("res_id = ?", routerDetail.UpstreamResID).
			Delete(&upstreamModel).Error; err != nil {
			return
		}

		upstreamNodeModel := models.UpstreamNodes{}
		if err = tx.Table(upstreamNodeModel.TableName()).
			Where("upstream_res_id = ?", routerDetail.UpstreamResID).
			Delete(&upstreamNodeModel).Error; err != nil {
			return
		}

		return
	})

	err = RouterRelease([]string{routerResId}, utils.ReleaseTypeDelete)
	if err != nil {
		return
	}

	err = UpstreamRelease([]string{routerDetail.UpstreamResID}, utils.ReleaseTypePush)
	if err != nil {
		return
	}

	return
}

func RouterCopy(routerResId string) (err error) {
	routerModel := models.Routers{}
	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	upstreamModel := models.Upstreams{}
	var upstreamDetail models.Upstreams
	if routerDetail.UpstreamResID != "" {
		upstreamDetail, err = upstreamModel.UpstreamDetailByResId(routerDetail.UpstreamResID)
		if err != nil {
			return
		}
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList := make([]models.UpstreamNodes, 0)
	if upstreamDetail.ResID != "" {
		upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamDetail.ResID})
		if err != nil {
			return
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList := make([]models.PluginConfigs, 0)
	pluginConfigList, err = pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, []string{routerResId})
	if err != nil {
		return
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {
		var upstreamResId string
		if upstreamDetail.ResID != "" {
			upstreamResId, err = upstreamModel.ModelUniqueId()
			if err != nil {
				return
			}

			err = tx.Table(upstreamModel.TableName()).Create(&models.Upstreams{
				ResID:          upstreamResId,
				Name:           upstreamDetail.Name,
				Algorithm:      upstreamDetail.Algorithm,
				ConnectTimeout: upstreamDetail.ConnectTimeout,
				WriteTimeout:   upstreamDetail.WriteTimeout,
				ReadTimeout:    upstreamDetail.ReadTimeout,
			}).Error
			if err != nil {
				return
			}

			newUpstreamNodeList := make([]models.UpstreamNodes, 0)
			if len(upstreamNodeList) != 0 {
				for _, upstreamNodeDetail := range upstreamNodeList {
					var upstreamNodeResId string
					upstreamNodeResId, err = upstreamNodeModel.ModelUniqueId()
					if err != nil {
						return
					}

					newUpstreamNodeList = append(newUpstreamNodeList, models.UpstreamNodes{
						ResID:         upstreamNodeResId,
						UpstreamResID: upstreamResId,
						NodeIP:        upstreamNodeDetail.NodeIP,
						IPType:        upstreamNodeDetail.IPType,
						NodePort:      upstreamNodeDetail.NodePort,
						NodeWeight:    upstreamNodeDetail.NodeWeight,
						Health:        upstreamNodeDetail.Health,
						HealthCheck:   upstreamNodeDetail.HealthCheck,
					})
				}

				if len(newUpstreamNodeList) != 0 {
					err = tx.Table(upstreamNodeModel.TableName()).Create(&newUpstreamNodeList).Error
					if err != nil {
						return
					}
				}
			}

		}

		newRouterResId, err := routerModel.ModelUniqueId()
		if err != nil {
			return
		}

		randomStr := utils.RandomStrGenerate(4)
		err = tx.Table(routerModel.TableName()).Create(&models.Routers{
			ResID:          newRouterResId,
			ServiceResID:   routerDetail.ServiceResID,
			UpstreamResID:  upstreamResId,
			RequestMethods: routerDetail.RequestMethods,
			RouterName:     routerDetail.RouterName + "-copy-" + randomStr,
			RouterPath:     routerDetail.RouterPath + "-copy-" + randomStr,
			Enable:         routerDetail.Enable,
			Release:        utils.ReleaseStatusU,
		}).Error
		if err != nil {
			return
		}

		newRouterPluginConfig := make([]models.PluginConfigs, 0)
		if len(pluginConfigList) > 0 {
			for _, pluginConfigInfo := range pluginConfigList {
				var pluginConfigresId string
				pluginConfigresId, err = pluginConfigModel.ModelUniqueId()
				if err != nil {
					return
				}

				newRouterPluginConfig= append(newRouterPluginConfig, models.PluginConfigs{
					ResID: pluginConfigresId,
					Name: pluginConfigInfo.Name,
					Type: models.PluginConfigsTypeRouter,
					TargetID: newRouterResId,
					PluginResID: pluginConfigInfo.PluginResID,
					PluginKey: pluginConfigInfo.PluginKey,
					Config: pluginConfigInfo.Config,
					Enable: pluginConfigInfo.Enable,
				})
			}

			err = tx.Table(pluginConfigModel.TableName()).Create(&newRouterPluginConfig).Error
			if err != nil {
				return
			}
		}

		return
	})

	return
}

func RouterInfoByResId(resId string) (models.Routers, error) {
	return (&models.Routers{}).RouterDetailByResId(resId)
}
