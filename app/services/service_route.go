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
	if routeInfo.ID != routeId {
		return errors.New(enums.CodeMessages(enums.RouteNull))
	}

	return nil
}

func CheckRouteEnableChange(routeId string, enable int) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.IsEnable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckRouteDelete(routeId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")

	if routeInfo.ReleaseStatus == utils.ReleaseStatusY {
		if routeInfo.IsEnable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if routeInfo.ReleaseStatus == utils.ReleaseStatusT {
		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
	}

	return nil
}

func CheckRouteRelease(routeId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.ReleaseStatus == utils.ReleaseStatusY {
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

func CheckExistServiceRoutePath(serviceId string, path string, filterRouteIds []string) error {
	routeModel := models.Routes{}
	routePaths, err := routeModel.RouteInfosByServiceRoutePath(serviceId, []string{path}, filterRouteIds)
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
		ServiceID:      routeData.ServiceID,
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		IsEnable:       routeData.IsEnable,
		ReleaseStatus:  utils.ReleaseStatusU,
	}

	if routeData.IsRelease == utils.IsReleaseY {
		createRouteData.ReleaseStatus = utils.ReleaseStatusY
	}

	routeId, err := createRouteData.RouteAdd(createRouteData)
	if err != nil {
		return err
	}

	if routeData.IsRelease == utils.IsReleaseY {
		configReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		if configReleaseErr != nil {
			createRouteData.ReleaseStatus = utils.ReleaseStatusU
			createRouteData.RouteUpdate(routeId, createRouteData)

			return configReleaseErr
		}
	}

	return nil
}

func RouteCopy(routeData *validators.ValidatorRouteAddUpdate, sourceRouteId string) error {
	routePluginModel := models.RoutePlugins{}
	routePluginInfos := routePluginModel.RoutePluginInfosByRouteId(sourceRouteId)

	createRouteData := models.Routes{
		ServiceID:      routeData.ServiceID,
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		IsEnable:       routeData.IsEnable,
		ReleaseStatus:  utils.ReleaseStatusU,
	}
	if routeData.IsRelease == utils.IsReleaseY {
		createRouteData.ReleaseStatus = utils.ReleaseStatusY
	}

	routeId, err := createRouteData.RouteCopy(createRouteData, routePluginInfos)
	if err != nil {
		return err
	}

	if routeData.IsRelease == utils.IsReleaseY {
		routeReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		if routeReleaseErr != nil {
			routeModel := models.Routes{}
			routeModel.ReleaseStatus = utils.ReleaseStatusU
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
		IsEnable:       routeData.IsEnable,
	}
	if len(routeData.RouteName) != 0 {
		updateRouteData.RouteName = routeData.RouteName
	}
	if routeInfo.ReleaseStatus == utils.ReleaseStatusY {
		updateRouteData.ReleaseStatus = utils.ReleaseStatusT
	}

	if routeData.IsRelease == utils.IsReleaseY {
		updateRouteData.ReleaseStatus = utils.ReleaseStatusY
	}

	err := routeModel.RouteUpdate(routeId, updateRouteData)
	if err != nil {
		return err
	}

	if routeData.IsRelease == utils.IsReleaseY {
		configReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, routeId)
		if configReleaseErr != nil {
			if routeInfo.ReleaseStatus != utils.ReleaseStatusU {
				updateRouteData.ReleaseStatus = utils.ReleaseStatusT
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

	routeInfo.ID = routeModelInfo.ID
	routeInfo.ServiceID = routeModelInfo.ServiceID
	routeInfo.RouteName = routeModelInfo.RouteName
	routeInfo.RequestMethods = strings.Split(routeModelInfo.RequestMethods, ",")
	routeInfo.RoutePath = routeModelInfo.RoutePath
	routeInfo.IsEnable = routeModelInfo.IsEnable
	routeInfo.ReleaseStatus = routeModelInfo.ReleaseStatus

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
		routeModel.RouteSwitchRelease(id, routeInfo.ReleaseStatus)
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

	routeConfig.ID = routeInfo.ID
	routeConfig.ServiceID = routeInfo.ServiceID
	routeConfig.Path = routeInfo.RoutePath
	routeConfig.IsEnable = routeInfo.IsEnable
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
		_, routePluginExist := routePluginAllPluginIdsMap[allPluginInfo.ID]

		if !routePluginExist {
			routeAddPluginInfo := RouteAddPluginInfo{}
			routeAddPluginInfo.ID = allPluginInfo.ID
			routeAddPluginInfo.Name = allPluginInfo.Name
			routeAddPluginInfo.Tag = allPluginInfo.Tag
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
		routePluginInfo.PluginId = routePluginConfigInfo.Plugin.ID
		routePluginInfo.Name = routePluginConfigInfo.Plugin.Name
		routePluginInfo.Tag = routePluginConfigInfo.Plugin.Tag
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
