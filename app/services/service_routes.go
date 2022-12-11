package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func CheckRouteExist(routeId string, serviceId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, serviceId)
	if routeInfo.ResID != routeId {
		return errors.New(enums.CodeMessages(enums.RouteNull))
	}

	return nil
}

func CheckRouteEnableChange(routeId string, enable int) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckRouteDelete(routeId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")

	if routeInfo.Release == utils.ReleaseStatusY {
		if routeInfo.Enable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if routeInfo.Release == utils.ReleaseStatusT {
		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
	}

	return nil
}

func CheckRouteRelease(routeId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.Release == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func CheckServiceRoutePath(path string) error {
	if path == utils.DefaultRoutePath {
		return errors.New(enums.CodeMessages(enums.RouteDefaultPathNoPermission))
	}

	if strings.Index(path, utils.DefaultRoutePath) == 0 {
		return errors.New(enums.CodeMessages(enums.RouteDefaultPathForbiddenPrefix))
	}

	return nil
}

func CheckEditDefaultPathRoute(routeId string) error {
	routeModel := models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.RoutePath == utils.DefaultRoutePath {
		return errors.New(enums.CodeMessages(enums.RouteDefaultPathNoPermission))
	}

	return nil
}

func CheckExistServiceRoutePath(serviceResId string, path string, filterRouteResIds []string) error {
	routeModel := models.Routes{}
	routePaths, err := routeModel.RouteInfosByServiceRoutePath(serviceResId, []string{path}, filterRouteResIds)
	if err != nil {
		return err
	}

	if len(routePaths) == 0 {
		return nil
	}

	existRoutePath := make([]string, 0)
	tmpExistRoutePathMap := make(map[string]byte, 0)
	for _, routePath := range routePaths {
		_, exist := tmpExistRoutePathMap[routePath.RoutePath]
		if exist {
			continue
		}

		existRoutePath = append(existRoutePath, routePath.RoutePath)
		tmpExistRoutePathMap[routePath.RoutePath] = 0
	}

	if len(existRoutePath) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.RoutePathExist), strings.Join(existRoutePath, ",")))
	}

	return nil
}

func RouteCreate(routeData *validators.ValidatorRouteAddUpdate) error {
	createRouteData := models.Routes{
		ServiceResID:   routeData.ServiceResID,
		UpstreamResID:  routeData.UpstreamResID,
		RouteName:      routeData.RouteName,
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		Enable:         routeData.Enable,
		Release:        utils.ReleaseStatusU,
	}

	if routeData.Release == utils.ReleaseY {
		createRouteData.Release = utils.ReleaseStatusY
	}

	createUpstreamData := models.Upstreams{
		Algorithm:      routeData.LoadBalance,
		ConnectTimeout: routeData.ConnectTimeout,
		WriteTimeout:   routeData.WriteTimeout,
		ReadTimeout:    routeData.ReadTimeout,
	}

	createUpstreamNodes := make([]models.UpstreamNodes, 0)
	if len(routeData.UpstreamNodes) > 0 {
		for _, upstreamNode := range routeData.UpstreamNodes {
			ipType, err := utils.DiscernIP(upstreamNode.NodeIp)
			if err != nil {
				return err
			}
			ipTypeMap := models.IPTypeMap()

			createUpstreamNodes = append(createUpstreamNodes, models.UpstreamNodes{
				NodeIP:     upstreamNode.NodeIp,
				IPType:     ipTypeMap[ipType],
				NodePort:   upstreamNode.NodePort,
				NodeWeight: upstreamNode.NodeWeight,
			})
		}
	}

	_, err := createRouteData.RouteAdd(createRouteData, createUpstreamData, createUpstreamNodes)

	if err != nil {
		return err
	}

	// @todo如果设置了发布，则这里需要发布 routeResId
	if routeData.Release == utils.ReleaseY {

	}

	return nil
}

func RouteCopy(routeData *validators.ValidatorRouteAddUpdate, sourceRouteId string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfos := routePluginModel.RoutePluginInfosByRouteId(sourceRouteId)

	createRouteData := models.Routes{
		ServiceResID:   routeData.ServiceResID,
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		Enable:         routeData.Enable,
		Release:        utils.ReleaseStatusU,
	}
	if routeData.Release == utils.ReleaseY {
		createRouteData.Release = utils.ReleaseStatusY
	}

	routeId, err := createRouteData.RouteCopy(createRouteData, routePluginInfos)
	if err != nil {
		return err
	}

	if routeData.Release == utils.ReleaseY {
		routeReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		if routeReleaseErr != nil {
			routeModel := models.Routes{}
			routeModel.Release = utils.ReleaseStatusU
			routeUpdateErr := routeModel.RouteUpdate(routeId, routeModel)
			if routeUpdateErr != nil {
				return routeUpdateErr
			}
		}
		return routeReleaseErr
	}

	return nil
}

func RouteUpdate(routeId string, routeData *validators.ValidatorRouteAddUpdate) error {
	routeModel := models.Routes{}
	routeInfo, routeInfoErr := routeModel.RouteInfoById(routeId)
	if routeInfoErr != nil {
		return routeInfoErr
	}

	updateRouteData := models.Routes{
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		Enable:         routeData.Enable,
	}
	if len(routeData.RouteName) != 0 {
		updateRouteData.RouteName = routeData.RouteName
	}
	if routeInfo.Release == utils.ReleaseStatusY {
		updateRouteData.Release = utils.ReleaseStatusT
	}

	if routeData.Release == utils.ReleaseY {
		updateRouteData.Release = utils.ReleaseStatusY
	}

	err := routeModel.RouteUpdate(routeId, updateRouteData)
	if err != nil {
		return err
	}

	if routeData.Release == utils.ReleaseY {
		configReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		if configReleaseErr != nil {
			if routeInfo.Release != utils.ReleaseStatusU {
				updateRouteData.Release = utils.ReleaseStatusT
			}
			routeModel.RouteUpdate(routeId, updateRouteData)

			return configReleaseErr
		}
	}

	return nil
}

func RouteDelete(routeId string) error {
	configReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypeDelete, routeId)
	if configReleaseErr != nil {
		return configReleaseErr
	}

	routeModel := models.Routes{}
	err := routeModel.RouteDelete(routeId)
	if err != nil {
		ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		return err
	}

	return nil
}

type routePlugin struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Tag           string `json:"tag"`
	Type          int    `json:"type"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
}

type StructRouteList struct {
	ID             string        `json:"id"`
	RouteName      string        `json:"route_name"`
	RequestMethods []string      `json:"request_methods"`
	RoutePath      string        `json:"route_path"`
	IsEnable       int           `json:"is_enable"`
	ReleaseStatus  int           `json:"release_status"`
	PluginList     []routePlugin `json:"plugin_list"`
}

func (s *StructRouteList) RouteListPage(serviceId string, param *validators.ValidatorRouteList) ([]StructRouteList, int, error) {

	routeModel := models.Routes{}
	routeInfos, total, listError := routeModel.RouteListPage(serviceId, param)

	routeList := make([]StructRouteList, 0)
	if len(routeInfos) != 0 {
		for _, routeInfo := range routeInfos {
			structRouteList := StructRouteList{}
			structRouteList.ID = routeInfo.ID
			structRouteList.RouteName = routeInfo.RouteName
			structRouteList.RequestMethods = strings.Split(routeInfo.RequestMethods, ",")
			structRouteList.RoutePath = routeInfo.RoutePath
			structRouteList.IsEnable = routeInfo.IsEnable
			structRouteList.ReleaseStatus = routeInfo.ReleaseStatus

			routePluginInfos := make([]routePlugin, 0)
			if len(routeInfo.Plugins) != 0 {
				for _, routePluginInfo := range routeInfo.Plugins {
					tmpRoutePluginInfo := routePlugin{}
					tmpRoutePluginInfo.ID = routePluginInfo.ID
					tmpRoutePluginInfo.Name = routePluginInfo.Name
					tmpRoutePluginInfo.Icon = routePluginInfo.Icon
					tmpRoutePluginInfo.Tag = routePluginInfo.Tag
					tmpRoutePluginInfo.Type = routePluginInfo.Type
					tmpRoutePluginInfo.IsEnable = routePluginInfo.IsEnable
					tmpRoutePluginInfo.ReleaseStatus = routePluginInfo.ReleaseStatus
					routePluginInfos = append(routePluginInfos, tmpRoutePluginInfo)
				}
			}
			structRouteList.PluginList = routePluginInfos

			routeList = append(routeList, structRouteList)
		}
	}

	return routeList, total, listError
}

type StructRouteInfo struct {
	ID             string   `json:"id"`
	ServiceID      string   `json:"service_id"`
	RouteName      string   `json:"route_name"`
	RequestMethods []string `json:"request_methods"`
	RoutePath      string   `json:"route_path"`
	IsEnable       int      `json:"is_enable"`
	ReleaseStatus  int      `json:"release_status"`
}

func (s *StructRouteInfo) RouteInfoByServiceRouteId(serviceId string, routeId string) (StructRouteInfo, error) {
	routeInfo := StructRouteInfo{}
	routeModel := &models.Routes{}
	routeModelInfo, routeModelInfoErr := routeModel.RouteInfosByServiceRouteId(serviceId, routeId)
	if routeModelInfoErr != nil {
		return routeInfo, routeModelInfoErr
	}

	routeInfo.ID = routeModelInfo.ResID
	routeInfo.ServiceID = routeModelInfo.ServiceResID
	routeInfo.RouteName = routeModelInfo.RouteName
	routeInfo.RequestMethods = strings.Split(routeModelInfo.RequestMethods, ",")
	routeInfo.RoutePath = routeModelInfo.RoutePath
	routeInfo.IsEnable = routeModelInfo.Enable
	routeInfo.ReleaseStatus = routeModelInfo.Release

	return routeInfo, nil
}

func ServiceRouteRelease(id string) error {
	routeModel := models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(id, "")
	err := routeModel.RouteSwitchRelease(id, utils.ReleaseStatusY)
	if err != nil {
		return err
	}

	configReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, id)
	if configReleaseErr != nil {
		routeModel.RouteSwitchRelease(id, routeInfo.Release)
		return configReleaseErr
	}

	return nil
}

func ServiceRouteConfigRelease(releaseType string, id string) error {
	routeConfig, routeConfigErr := generateRouteConfig(id)
	if routeConfigErr != nil {
		return routeConfigErr
	}

	routeConfigJson, routeConfigJsonErr := json.Marshal(routeConfig)
	if routeConfigJsonErr != nil {
		return routeConfigJsonErr
	}
	routeConfigStr := string(routeConfigJson)

	etcdKey := utils.EtcdKey(utils.EtcdKeyTypeRoute, id)
	if len(etcdKey) == 0 {
		return errors.New(enums.CodeMessages(enums.EtcdKeyNull))
	}

	etcdClient := packages.GetEtcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	defer cancel()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(ctx, etcdKey, routeConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(ctx, etcdKey)
	}

	if respErr != nil {
		return errors.New(enums.CodeMessages(enums.EtcdUnavailable))
	}

	return nil
}

type RouteConfig struct {
	ID        string   `json:"id"`
	ServiceID string   `json:"service_id"`
	Path      string   `json:"path"`
	IsEnable  int      `json:"is_enable"`
	Methods   []string `json:"methods"`
}

func generateRouteConfig(id string) (RouteConfig, error) {
	routeConfig := RouteConfig{}
	routeModel := models.Routes{}
	routeInfo, routeInfoErr := routeModel.RouteInfoById(id)
	if routeInfoErr != nil {
		return routeConfig, routeInfoErr
	}

	methods := strings.Split(routeInfo.RequestMethods, ",")

	var allMethod bool
	if len(methods) != 0 {
		for _, method := range methods {
			if method == utils.RequestMethodALL {
				allMethod = true
				break
			}
		}
	}

	if allMethod == true {
		methods = utils.ConfigAllRequestMethod()
	}

	routeConfig.ID = routeInfo.ResID
	routeConfig.ServiceID = routeInfo.ServiceResID
	routeConfig.Path = routeInfo.RoutePath
	routeConfig.IsEnable = routeInfo.Enable
	routeConfig.Methods = methods

	return routeConfig, nil
}

type RouteAddPluginInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

func (r *RouteAddPluginInfo) RouteAddPluginList(filterRouteId string) ([]RouteAddPluginInfo, error) {
	routeAddPluginList := make([]RouteAddPluginInfo, 0)
	if len(filterRouteId) == 0 {
		return routeAddPluginList, errors.New(enums.CodeMessages(enums.ParamsError))
	}

	pluginsModel := models.Plugins{}
	allPluginList := pluginsModel.PluginAllList()

	routePluginsModel := models.RoutePlugins{}
	routePluginAllList := routePluginsModel.RoutePluginAllListByRouteIds([]string{filterRouteId})

	routePluginAllPluginIdsMap := make(map[string]byte, 0)
	for _, routePluginInfo := range routePluginAllList {
		routePluginAllPluginIdsMap[routePluginInfo.PluginID] = 0
	}

	for _, allPluginInfo := range allPluginList {
		_, routePluginExist := routePluginAllPluginIdsMap[allPluginInfo.ResID]

		if !routePluginExist {
			routeAddPluginInfo := RouteAddPluginInfo{}
			routeAddPluginInfo.ID = allPluginInfo.ResID
			routeAddPluginInfo.Tag = allPluginInfo.PluginKey
			routeAddPluginInfo.Type = allPluginInfo.Type
			routeAddPluginInfo.Description = allPluginInfo.Description

			routeAddPluginList = append(routeAddPluginList, routeAddPluginInfo)
		}
	}

	return routeAddPluginList, nil
}

type RoutePluginInfo struct {
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

func (r *RoutePluginInfo) RoutePluginList(routeId string) []RoutePluginInfo {
	routePluginList := make([]RoutePluginInfo, 0)

	routePluginsModel := models.RoutePlugins{}
	routePluginConfigInfos := routePluginsModel.RoutePluginInfoConfigListByRouteIds([]string{routeId})

	for _, routePluginConfigInfo := range routePluginConfigInfos {
		routePluginInfo := RoutePluginInfo{}
		routePluginInfo.ID = routePluginConfigInfo.ID
		routePluginInfo.PluginId = routePluginConfigInfo.Plugin.ResID
		routePluginInfo.Tag = routePluginConfigInfo.Plugin.PluginKey
		routePluginInfo.Icon = routePluginConfigInfo.Plugin.Icon
		routePluginInfo.Type = routePluginConfigInfo.Plugin.Type
		routePluginInfo.Description = routePluginConfigInfo.Plugin.Description
		routePluginInfo.Order = routePluginConfigInfo.Order
		routePluginInfo.Config = routePluginConfigInfo.Config
		routePluginInfo.IsEnable = routePluginConfigInfo.IsEnable
		routePluginInfo.ReleaseStatus = routePluginConfigInfo.ReleaseStatus

		routePluginList = append(routePluginList, routePluginInfo)
	}

	return routePluginList
}
