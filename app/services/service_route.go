package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
	"fmt"
	"strings"
)

func CheckExistRoutePath(path string, filterRouteIds []string) error {
	routeModel := models.Routes{}
	routePaths, err := routeModel.RouteInfosByRoutePath([]string{path}, filterRouteIds)
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

func RouteCreate(routeData *validators.RouteAddUpdate) error {

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
