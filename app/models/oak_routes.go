package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
)

type Routes struct {
	ID             string `gorm:"column:id;primary_key"`  //Route id
	ServiceID      string `gorm:"column:service_id"`      //Service id
	RouteName      string `gorm:"column:route_name"`      //Route name
	RequestMethods string `gorm:"column:request_methods"` //Request method
	RoutePath      string `gorm:"column:route_path"`      //Routing path
	IsEnable       int    `gorm:"column:is_enable"`       //Routing enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (r *Routes) TableName() string {
	return "oak_routes"
}

var routeId = ""

func (r *Routes) RouteIdUnique(routeIds map[string]string) (string, error) {
	if r.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeRoute)
		if err != nil {
			return "", err
		}
		r.ID = tmpID
	}

	result := packages.GetDb().Table(r.TableName()).Select("id").First(&r)
	mapId := routeIds[r.ID]
	if (result.RowsAffected == 0) && (r.ID != mapId) {
		routeId = r.ID
		routeIds[r.ID] = r.ID
		return routeId, nil
	} else {
		svcId, svcErr := utils.IdGenerate(utils.IdTypeService)
		if svcErr != nil {
			return "", svcErr
		}
		r.ID = svcId
		_, err := r.RouteIdUnique(routeIds)
		if err != nil {
			return "", err
		}
	}

	return routeId, nil
}

func (r *Routes) RouteInfosByRoutePath(routePaths []string, filterRouteIds []string) ([]Routes, error) {
	routesInfos := make([]Routes, 0)

	db := packages.GetDb().Table(r.TableName()).Where("route_path IN ?", routePaths)
	if len(filterRouteIds) != 0 {
		db = db.Where("id NOT IN ?", filterRouteIds)
	}
	err := db.Find(&routesInfos).Error

	return routesInfos, err
}

func (r *Routes) RouteAdd(routeData Routes) error {
	tpmIds := map[string]string{}
	routeId, routeIdUniqueErr := r.RouteIdUnique(tpmIds)
	if routeIdUniqueErr != nil {
		return routeIdUniqueErr
	}

	routeData.ID = routeId
	routeData.RouteName = routeId

	addErr := packages.GetDb().Table(r.TableName()).Create(&routeData).Error
	if addErr != nil {
		return addErr
	}

	return nil
}
