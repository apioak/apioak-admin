package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"fmt"
	"strings"
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

func CheckRouteEnableOnProhibitedOp(routeId string) error {
	routeModel := &models.Routes{}
	routeInfo := routeModel.RouteInfoByIdServiceId(routeId, "")
	if routeInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
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
	}

	routeModel := models.Routes{}
	err := routeModel.RouteAdd(createRouteData)
	if err != nil {
		return err
	}

	// @todo 如果状态是"开启"，则需要同步远程数据中心

	return nil
}

func RouteUpdate(id string, routeData *validators.ValidatorRouteAddUpdate) error {
	updateRouteData := models.Routes{
		RequestMethods: routeData.RequestMethods,
		RoutePath:      routeData.RoutePath,
		IsEnable:       routeData.IsEnable,
	}
	if len(routeData.RouteName) != 0 {
		updateRouteData.RouteName = routeData.RouteName
	}

	routeModel := models.Routes{}
	err := routeModel.RouteUpdate(id, updateRouteData)
	if err != nil {
		return err
	}

	// @todo 如果状态是"开启"，则需要同步远程数据中心

	return nil
}

func RouteDelete(id string) error {
	routeModel := models.Routes{}
	err := routeModel.RouteDelete(id)
	if err != nil {
		return err
	}

	// @todo 如果状态是"开启"，则需要同步远程数据中心

	return nil
}

type routePlugin struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type StructRouteList struct {
	ID             string        `json:"id"`
	RouteName      string        `json:"route_name"`
	RequestMethods []string      `json:"request_methods"`
	RoutePath      string        `json:"route_path"`
	IsEnable       int           `json:"is_enable"`
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

			routePluginInfos := make([]routePlugin, 0)
			if len(routeInfo.Plugins) != 0 {
				for _, routePluginInfo := range routeInfo.Plugins {
					tmpRoutePluginInfo := routePlugin{}
					tmpRoutePluginInfo.ID = routePluginInfo.ID
					tmpRoutePluginInfo.Name = routePluginInfo.Name
					tmpRoutePluginInfo.Icon = routePluginInfo.Icon
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

	return routeInfo, nil
}

func ServiceRouteSwitchEnable(id string, enable int) error {
	routeModel := models.Routes{}
	err := routeModel.RouteSwitchEnable(id, enable)
	if err != nil {
		return err
	}

	// @todo 触发远程发布数据

	return nil
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Config      string `json:"config"`
	IsEnable    int    `json:"is_enable"`
}

func (r *RoutePluginInfo) RoutePluginList(routeId string) []RoutePluginInfo {
	routePluginList := make([]RoutePluginInfo, 0)

	routePluginsModel := models.RoutePlugins{}
	routePluginConfigInfos := routePluginsModel.RoutePluginInfoConfigListByRouteIds([]string{routeId})

	for _, routePluginConfigInfo := range routePluginConfigInfos {
		routePluginInfo := RoutePluginInfo{}
		routePluginInfo.ID = routePluginConfigInfo.Plugin.ID
		routePluginInfo.Name = routePluginConfigInfo.Plugin.Name
		routePluginInfo.Tag = routePluginConfigInfo.Plugin.Tag
		routePluginInfo.Type = routePluginConfigInfo.Plugin.Type
		routePluginInfo.Description = routePluginConfigInfo.Plugin.Description
		routePluginInfo.Order = routePluginConfigInfo.Order
		routePluginInfo.Config = routePluginConfigInfo.Config
		routePluginInfo.IsEnable = routePluginConfigInfo.IsEnable

		routePluginList = append(routePluginList, routePluginInfo)
	}

	return routePluginList
}
