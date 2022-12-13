package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func CheckRouterExist(routerResId string, serviceResId string) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterInfoByResIdServiceResId(routerResId, serviceResId)

	if len(routerInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterNull))
	}

	return nil
}

// func CheckRouterEnableChange(routerId string, enable int) error {
// 	routerModel := &models.Routers{}
// 	routerInfo := routerModel.RouterInfoByResIdServiceResId(routerId, "")
//
// 	if routerInfo.Enable == enable {
// 		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
// 	}
//
// 	return nil
// }
//
// func CheckRouterDelete(routerId string) error {
// 	routerModel := &models.Routers{}
// 	routerInfo := routerModel.RouterInfoByResIdServiceResId(routerId, "")
//
// 	if routerInfo.Release == utils.ReleaseStatusY {
// 		if routerInfo.Enable == utils.EnableOn {
// 			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
// 		}
// 	} else if routerInfo.Release == utils.ReleaseStatusT {
// 		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
// 	}
//
// 	return nil
// }

func CheckRouterRelease(routerResId string) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterInfoByResIdServiceResId(routerResId, "")

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

func CheckEditDefaultPathRouter(routerId string) error {
	routerModel := models.Routers{}
	routerInfo := routerModel.RouterInfoByResIdServiceResId(routerId, "")
	if routerInfo.RouterPath == utils.DefaultRouterPath {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathNoPermission))
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

func RouterCreate(routerData *validators.ValidatorRouterAddUpdate) error {
	createRouterData := models.Routers{
		ServiceResID:   routerData.ServiceResID,
		UpstreamResID:  routerData.UpstreamResID,
		RouterName:     routerData.RouterName,
		RequestMethods: routerData.RequestMethods,
		RouterPath:     routerData.RouterPath,
		Enable:         routerData.Enable,
		Release:        utils.ReleaseStatusU,
	}

	if routerData.Release == utils.ReleaseY {
		createRouterData.Release = utils.ReleaseStatusY
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
			ipType, err := utils.DiscernIP(upstreamNode.NodeIp)
			if err != nil {
				return err
			}
			ipTypeMap := models.IPTypeMap()

			createUpstreamNodes = append(createUpstreamNodes, models.UpstreamNodes{
				NodeIP:      upstreamNode.NodeIp,
				IPType:      ipTypeMap[ipType],
				NodePort:    upstreamNode.NodePort,
				NodeWeight:  upstreamNode.NodeWeight,
				Health:      upstreamNode.Health,
				HealthCheck: upstreamNode.HealthCheck,
			})
		}
	}

	_, err := createRouterData.RouterAdd(createRouterData, createUpstreamData, createUpstreamNodes)

	if err != nil {
		return err
	}

	return nil
}

// func RouterCopy(routerData *validators.ValidatorRouterAddUpdate, sourceRouterId string) error {
// 	routerPluginModel := models.RouterPlugins{}
// 	routerPluginInfos := routerPluginModel.RouterPluginInfosByRouterId(sourceRouterId)
//
// 	createRouterData := models.Routers{
// 		ServiceResID:   routerData.ServiceResID,
// 		RequestMethods: routerData.RequestMethods,
// 		RouterPath:      routerData.RouterPath,
// 		Enable:         routerData.Enable,
// 		Release:        utils.ReleaseStatusU,
// 	}
// 	if routerData.Release == utils.ReleaseY {
// 		createRouterData.Release = utils.ReleaseStatusY
// 	}
//
// 	routerId, err := createRouterData.RouterCopy(createRouterData, routerPluginInfos)
// 	if err != nil {
// 		return err
// 	}
//
// 	if routerData.Release == utils.ReleaseY {
// 		routerReleaseErr := ServiceRouterConfigRelease(utils.ReleaseTypePush, routerId)
// 		if routerReleaseErr != nil {
// 			routerModel := models.Routers{}
// 			routerModel.Release = utils.ReleaseStatusU
// 			routerUpdateErr := routerModel.RouterUpdate(routerId, routerModel)
// 			if routerUpdateErr != nil {
// 				return routerUpdateErr
// 			}
// 		}
// 		return routerReleaseErr
// 	}
//
// 	return nil
// }

// func RouterUpdate(routerId string, routerData *validators.ValidatorRouterAddUpdate) error {
// 	routerModel := models.Routers{}
// 	routerInfo, routerInfoErr := routerModel.RouterInfoById(routerId)
// 	if routerInfoErr != nil {
// 		return routerInfoErr
// 	}
//
// 	updateRouterData := models.Routers{
// 		RequestMethods: routerData.RequestMethods,
// 		RouterPath:      routerData.RouterPath,
// 		Enable:         routerData.Enable,
// 	}
// 	if len(routerData.RouterName) != 0 {
// 		updateRouterData.RouterName = routerData.RouterName
// 	}
// 	if routerInfo.Release == utils.ReleaseStatusY {
// 		updateRouterData.Release = utils.ReleaseStatusT
// 	}
//
// 	if routerData.Release == utils.ReleaseY {
// 		updateRouterData.Release = utils.ReleaseStatusY
// 	}
//
// 	err := routerModel.RouterUpdate(routerId, updateRouterData)
// 	if err != nil {
// 		return err
// 	}
//
// 	if routerData.Release == utils.ReleaseY {
// 		configReleaseErr := ServiceRouterConfigRelease(utils.ReleaseTypePush, routerId)
// 		if configReleaseErr != nil {
// 			if routerInfo.Release != utils.ReleaseStatusU {
// 				updateRouterData.Release = utils.ReleaseStatusT
// 			}
// 			routerModel.RouterUpdate(routerId, updateRouterData)
//
// 			return configReleaseErr
// 		}
// 	}
//
// 	return nil
// }
//
// func RouterDelete(routerId string) error {
// 	configReleaseErr := ServiceRouterConfigRelease(utils.ReleaseTypeDelete, routerId)
// 	if configReleaseErr != nil {
// 		return configReleaseErr
// 	}
//
// 	routerModel := models.Routers{}
// 	err := routerModel.RouterDelete(routerId)
// 	if err != nil {
// 		ServiceRouterConfigRelease(utils.ReleaseTypePush, routerId)
// 		return err
// 	}
//
// 	return nil
// }

type routerPlugin struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Tag           string `json:"tag"`
	Type          int    `json:"type"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
}

type RouterListItem struct {
	ResID          string         `json:"res_id"`
	ServiceResId   string         `json:"service_res_id"`
	ServiceName    string         `json:"service_name"`
	RouterName     string         `json:"router_name"`
	RequestMethods []string       `json:"request_methods"`
	RouterPath     string         `json:"router_path"`
	Enable         int            `json:"enable"`
	Release        int            `json:"release"`
	PluginList     []routerPlugin `json:"plugin_list"`
}

func (s *RouterListItem) RouterListPage(serviceResId string, param *validators.ValidatorRouterList) ([]RouterListItem, int, error) {

	routerModel := models.Routers{}
	routerInfos, total, listError := routerModel.RouterListPage(serviceResId, param)

	routerList := make([]RouterListItem, 0)
	if len(routerInfos) != 0 {
		for _, routerInfo := range routerInfos {
			routerListItem := RouterListItem{}
			routerListItem.ResID = routerInfo.ResID
			routerListItem.RouterName = routerInfo.RouterName
			routerListItem.RequestMethods = strings.Split(routerInfo.RequestMethods, ",")
			routerListItem.RouterPath = routerInfo.RouterPath
			routerListItem.Enable = routerInfo.Enable
			routerListItem.Release = routerInfo.Release

			// @todo 这里补充路由的插件列表数据，还有一个是服务的名称也需要补充
			// routerListItem.PluginList = routerPluginInfos

			routerList = append(routerList, routerListItem)
		}
	}

	return routerList, total, listError
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

func ServiceRouterRelease(routerResIds []string, releaseType string) error {
	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		return errors.New(enums.CodeMessages(enums.ReleaseTypeError))
	}

	routerModel := models.Routers{}
	routerList, err := routerModel.RouterListByRouterResIds(routerResIds)
	if err != nil {
		return err
	}

	if len(routerList) == 0 {
		return nil
	}

	serviceResIds := make([]string, 0)

	for _, routerInfo := range routerList {
		if len(routerInfo.ServiceResID) > 0 {
			serviceResIds = append(serviceResIds, routerInfo.ServiceResID)
		}
	}

	publishedServiceResIdsMap := make(map[string]byte)
	// @todo 根据服务ID获取已经发布的服务数据（如果没有已经发布的数据，则本次发布不允许，直接返回错误信息即可）

	toBeOpUpstreamResIds := make([]string, 0)
	toBeOpRouterList := make([]models.Routers, 0)
	for _, routerInfo := range routerList {

		_, ok := publishedServiceResIdsMap[routerInfo.ServiceResID]
		if !ok {
			continue
		}

		toBeOpRouterList = append(toBeOpRouterList, routerInfo)

		if len(routerInfo.UpstreamResID) > 0 {
			toBeOpUpstreamResIds = append(toBeOpUpstreamResIds, routerInfo.UpstreamResID)
		}
	}

	if len(toBeOpRouterList) == 0 {
		return nil
	}

	routerConfigList := make([]rpc.RouterConfig, 0)
	for _, toBeOpRouterInfo := range toBeOpRouterList {
		routerConfig, routerConfigErr := generateRouterConfig(toBeOpRouterInfo)
		if routerConfigErr != nil {
			return routerConfigErr
		}

		if len(routerConfig.Name) == 0 {
			continue
		}

		routerConfigList = append(routerConfigList, routerConfig)
	}

	newApiOak := rpc.NewApiOak()

	if releaseType == utils.ReleaseTypePush {
		releaseUpstreamErr := RouterUpstreamRelease(toBeOpUpstreamResIds, releaseType)
		if releaseUpstreamErr != nil {
			return releaseUpstreamErr
		}

		routerPutErr := newApiOak.RouterPut(routerConfigList)
		if routerPutErr != nil {
			return routerPutErr
		}

		for _, toBeOpRouterInfo := range toBeOpRouterList {
			switchReleaseErr := routerModel.RouterSwitchRelease(toBeOpRouterInfo.ResID, utils.ReleaseStatusY)
			if switchReleaseErr != nil {
				return switchReleaseErr
			}
		}

	} else {
		routerDeleteErr := newApiOak.RouterDelete(routerConfigList)
		if routerDeleteErr != nil {
			return routerDeleteErr
		}

		releaseUpstreamErr := RouterUpstreamRelease(toBeOpUpstreamResIds, releaseType)
		if releaseUpstreamErr != nil {
			return releaseUpstreamErr
		}
	}

	return nil
}

func ServiceRouterConfigRelease(releaseType string, id string) error {

	// routerConfig, routerConfigErr := generateRouterConfig(id)
	// if routerConfigErr != nil {
	// 	return routerConfigErr
	// }
	routerConfig := rpc.RouterConfig{}

	routerConfigJson, routerConfigJsonErr := json.Marshal(routerConfig)
	if routerConfigJsonErr != nil {
		return routerConfigJsonErr
	}
	routerConfigStr := string(routerConfigJson)

	etcdKey := utils.EtcdKey(utils.EtcdKeyTypeRouter, id)
	if len(etcdKey) == 0 {
		return errors.New(enums.CodeMessages(enums.EtcdKeyNull))
	}

	etcdClient := packages.GetEtcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	defer cancel()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(ctx, etcdKey, routerConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(ctx, etcdKey)
	}

	if respErr != nil {
		return errors.New(enums.CodeMessages(enums.EtcdUnavailable))
	}

	return nil
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

	// @todo 根据路由res_id获取插件列表数据进行补充插件数据
	routerConfig.Plugins = make([]rpc.ConfigObjectName, 0)

	return routerConfig, nil
}

type RouterAddPluginInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

// func (r *RouterAddPluginInfo) RouterAddPluginList(filterRouterId string) ([]RouterAddPluginInfo, error) {
// 	routerAddPluginList := make([]RouterAddPluginInfo, 0)
// 	if len(filterRouterId) == 0 {
// 		return routerAddPluginList, errors.New(enums.CodeMessages(enums.ParamsError))
// 	}
//
// 	pluginsModel := models.Plugins{}
// 	allPluginList := pluginsModel.PluginAllList()
//
// 	routerPluginsModel := models.RouterPlugins{}
// 	routerPluginAllList := routerPluginsModel.RouterPluginAllListByRouterIds([]string{filterRouterId})
//
// 	routerPluginAllPluginIdsMap := make(map[string]byte, 0)
// 	for _, routerPluginInfo := range routerPluginAllList {
// 		routerPluginAllPluginIdsMap[routerPluginInfo.PluginID] = 0
// 	}
//
// 	for _, allPluginInfo := range allPluginList {
// 		_, routerPluginExist := routerPluginAllPluginIdsMap[allPluginInfo.ResID]
//
// 		if !routerPluginExist {
// 			routerAddPluginInfo := RouterAddPluginInfo{}
// 			routerAddPluginInfo.ID = allPluginInfo.ResID
// 			routerAddPluginInfo.Tag = allPluginInfo.PluginKey
// 			routerAddPluginInfo.Type = allPluginInfo.Type
// 			routerAddPluginInfo.Description = allPluginInfo.Description
//
// 			routerAddPluginList = append(routerAddPluginList, routerAddPluginInfo)
// 		}
// 	}
//
// 	return routerAddPluginList, nil
// }

type RouterPluginInfo struct {
	ID            string `json:"id"`
	PluginId      string `json:"plugin_id"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Icon          string `json:"icon"`
	Type          int    `json:"type"`
	Description   string `json:"description"`
	Order         int    `json:"order"`
	Config        string `json:"config"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
}

// func (r *RouterPluginInfo) RouterPluginList(routerId string) []RouterPluginInfo {
// 	routerPluginList := make([]RouterPluginInfo, 0)
//
// 	routerPluginsModel := models.RouterPlugins{}
// 	routerPluginConfigInfos := routerPluginsModel.RouterPluginInfoConfigListByRouterIds([]string{routerId})
//
// 	for _, routerPluginConfigInfo := range routerPluginConfigInfos {
// 		routerPluginInfo := RouterPluginInfo{}
// 		routerPluginInfo.ID = routerPluginConfigInfo.ID
// 		routerPluginInfo.PluginId = routerPluginConfigInfo.Plugin.ResID
// 		routerPluginInfo.Tag = routerPluginConfigInfo.Plugin.PluginKey
// 		routerPluginInfo.Icon = routerPluginConfigInfo.Plugin.Icon
// 		routerPluginInfo.Type = routerPluginConfigInfo.Plugin.Type
// 		routerPluginInfo.Description = routerPluginConfigInfo.Plugin.Description
// 		routerPluginInfo.Order = routerPluginConfigInfo.Order
// 		routerPluginInfo.Config = routerPluginConfigInfo.Config
// 		routerPluginInfo.IsEnable = routerPluginConfigInfo.IsEnable
// 		routerPluginInfo.ReleaseStatus = routerPluginConfigInfo.ReleaseStatus
//
// 		routerPluginList = append(routerPluginList, routerPluginInfo)
// 	}
//
// 	return routerPluginList
// }
