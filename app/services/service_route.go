package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
	"fmt"
	"strings"
)

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

func (s *StructRouteList) RouteListPage(
	serviceId string,
	param *validators.ValidatorRouteList) ([]StructRouteList, int, error) {

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
