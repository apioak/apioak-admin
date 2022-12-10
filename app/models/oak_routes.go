package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strings"
	"time"
)

type Routes struct {
	ID             string `gorm:"column:id;primary_key"`  //Route id
	ServiceID      string `gorm:"column:service_id"`      //Service id
	RouteName      string `gorm:"column:route_name"`      //Route name
	RequestMethods string `gorm:"column:request_methods"` //Request method
	RoutePath      string `gorm:"column:route_path"`      //Routing path
	IsEnable       int    `gorm:"column:is_enable"`       //Routing enable  1:on  2:off
	ReleaseStatus  int    `gorm:"column:release_status"`  //Route release status 1:unpublished  2:to be published  3:published
	ModelTime
}

type RoutesPlugins struct {
	Routes
	Plugins []Plugins `gorm:"many2many:oak_route_plugins;foreignKey:ID;joinForeignKey:RouteID;References:ID;JoinReferences:PluginID"`
}

// TableName sets the insert table name for this struct type
func (r *Routes) TableName() string {
	return "oak_routes"
}

var recursionTimesRoute = 1

func (m *Routes) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeRoute)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("id = ?", generateId).
		Select("id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesRoute = 1
		return generateId, nil
	} else {
		if recursionTimesRoute == utils.IdGenerateMaxTimes {
			recursionTimesRoute = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesRoute++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (r *Routes) RouteInfosByServiceRoutePath(
	serviceId string,
	routePaths []string,
	filterRouteIds []string) ([]Routes, error) {
	routesInfos := make([]Routes, 0)
	db := packages.GetDb().
		Table(r.TableName()).
		Where("service_id = ?", serviceId).
		Where("route_path IN ?", routePaths)

	if len(filterRouteIds) != 0 {
		db = db.Where("id NOT IN ?", filterRouteIds)
	}

	err := db.Find(&routesInfos).Error

	return routesInfos, err
}

func (r *Routes) RouteInfoById(routeId string) (Routes, error) {
	routeInfo := Routes{}
	err := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", routeId).
		First(&routeInfo).Error

	return routeInfo, err
}

func (r *Routes) RouteInfoByIdServiceId(routeId string, serviceId string) Routes {
	routeInfo := Routes{}
	db := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", routeId)

	if len(serviceId) != 0 {
		db = db.Where("service_id = ?", serviceId)
	}

	db.First(&routeInfo)

	return routeInfo
}

func (r *Routes) RouteInfosByServiceRouteId(serviceId string, routeId string) (Routes, error) {
	routeInfo := Routes{}
	err := packages.GetDb().
		Table(r.TableName()).
		Where("service_id = ?", serviceId).
		Where("id = ?", routeId).
		First(&routeInfo).Error

	return routeInfo, err
}

func (r *Routes) RouteInfosByServiceIdReleaseStatus(serviceId string, releaseStatus []int) []Routes {
	routeInfos := make([]Routes, 0)
	if len(serviceId) == 0 {
		return routeInfos
	}

	db := packages.GetDb().
		Table(r.TableName()).
		Where("service_id = ?", serviceId)

	if len(releaseStatus) != 0 {
		db = db.Where("release_status IN ?", releaseStatus)
	}
	db.Find(&routeInfos)

	return routeInfos
}

func (r *Routes) RouteAdd(routeData Routes) (string, error) {
	routeId, routeIdUniqueErr := r.ModelUniqueId()
	if routeIdUniqueErr != nil {
		return routeId, routeIdUniqueErr
	}

	routeData.ID = routeId
	routeData.RouteName = routeId

	addErr := packages.GetDb().
		Table(r.TableName()).
		Create(&routeData).Error

	if addErr != nil {
		return routeId, addErr
	}

	return routeId, nil
}

func (r *Routes) RouteUpdate(id string, routeData Routes) error {
	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(&routeData).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *Routes) RouteCopy(routeData Routes, routePlugins []RoutePlugins) (string, error) {
	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return "", err
	}

	routeId, routeIdUniqueErr := r.ModelUniqueId()
	if routeIdUniqueErr != nil {
		return routeId, routeIdUniqueErr
	}

	routeData.ID = routeId
	routeData.RouteName = routeId

	addErr := tx.
		Table(r.TableName()).
		Create(&routeData).Error

	if addErr != nil {
		tx.Rollback()
		return routeId, addErr
	}

	if len(routePlugins) != 0 {
		routePluginModel := RoutePlugins{}
		for k, _ := range routePlugins {
			routePluginId, routePluginIdErr := routePluginModel.ModelUniqueId()
			if routePluginIdErr != nil {
				tx.Rollback()
				return routeId, routePluginIdErr
			}

			routePlugins[k].ID = routePluginId
			routePlugins[k].RouteID = routeId
			routePlugins[k].ReleaseStatus = utils.ReleaseStatusU
			routePlugins[k].CreatedAt = time.Now()
			routePlugins[k].UpdatedAt = time.Now()
		}

		addRoutePluginErr := tx.
			Table(routePluginModel.TableName()).
			Create(&routePlugins).Error

		if addRoutePluginErr != nil {
			tx.Rollback()
			return routeId, addRoutePluginErr
		}
	}

	return routeId, tx.Commit().Error
}

func (r *Routes) RouteDelete(id string) error {

	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	deleteRouteError := tx.
		Table(r.TableName()).
		Where("id = ?", id).
		Delete(Routes{}).Error

	if deleteRouteError != nil {
		tx.Rollback()
		return deleteRouteError
	}

	routePluginsModel := RoutePlugins{}
	deleteRoutePluginError := tx.
		Table(routePluginsModel.TableName()).
		Where("route_id = ?", id).
		Delete(routePluginsModel).Error

	if deleteRoutePluginError != nil {
		tx.Rollback()
		return deleteRoutePluginError
	}

	return tx.Commit().Error
}

type RoutePluginConfigs struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Icon          string `json:"icon"`
	Type          int    `json:"type"`
	Config        string `json:"config"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
}

type RoutePluginConfigList struct {
	ID             string               `json:"id"`
	ServiceID      string               `json:"service_id"`
	RouteName      string               `json:"route_name"`
	RequestMethods string               `json:"request_methods"`
	RoutePath      string               `json:"route_path"`
	IsEnable       int                  `json:"is_enable"`
	ReleaseStatus  int                  `json:"release_status"`
	Plugins        []RoutePluginConfigs `json:"plugins"`
}

func (r *Routes) RouteListPage(serviceId string, param *validators.ValidatorRouteList) (list []RoutePluginConfigList, total int, listError error) {
	list = make([]RoutePluginConfigList, 0)
	routesPluginList := make([]RoutesPlugins, 0)

	routesPlugins := RoutesPlugins{}
	tx := packages.GetDb().
		Table(routesPlugins.TableName()).
		Where("service_id = ?", serviceId)

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		orWhere := packages.GetDb().
			Or("id LIKE ?", search).
			Or("route_name LIKE ?", search).
			Or("route_path LIKE ?", search).
			Or("request_methods LIKE ?", strings.ToUpper(search))

		tx = tx.Where(orWhere)
	}

	if param.IsEnable != 0 {
		tx = tx.Where("is_enable = ?", param.IsEnable)
	}
	if param.ReleaseStatus != 0 {
		tx = tx.Where("release_status = ?", param.ReleaseStatus)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.
		Preload("Plugins").
		Order("created_at desc")

	listError = ListPaginate(tx, &routesPluginList, &param.BaseListPage)

	if len(routesPluginList) == 0 {
		return
	}

	routeIds := make([]string, 0)
	for _, routesPluginInfo := range routesPluginList {
		routeIds = append(routeIds, routesPluginInfo.ID)
	}

	routePluginsModel := RoutePlugins{}
	routesPluginConfigList := routePluginsModel.RoutePluginAllListByRouteIds(routeIds)

	routesPluginConfigMap := make(map[string]map[string]RoutePlugins)
	if len(routesPluginConfigList) != 0 {
		for _, routesPluginConfigInfo := range routesPluginConfigList {
			if len(routesPluginConfigMap[routesPluginConfigInfo.RouteID]) == 0 {
				routesPluginConfigMap[routesPluginConfigInfo.RouteID] = make(map[string]RoutePlugins)
			}
			routesPluginConfigMap[routesPluginConfigInfo.RouteID][routesPluginConfigInfo.PluginID] = routesPluginConfigInfo
		}
	}

	for _, routesPluginInfo := range routesPluginList {
		routePluginConfigList := RoutePluginConfigList{
			ID:             routesPluginInfo.ID,
			ServiceID:      routesPluginInfo.ServiceID,
			RouteName:      routesPluginInfo.RouteName,
			RequestMethods: routesPluginInfo.RequestMethods,
			RoutePath:      routesPluginInfo.RoutePath,
			IsEnable:       routesPluginInfo.IsEnable,
			ReleaseStatus:  routesPluginInfo.ReleaseStatus,
		}

		routePluginsList := make([]RoutePluginConfigs, 0)
		if len(routesPluginInfo.Plugins) != 0 {
			for _, pluginInfo := range routesPluginInfo.Plugins {
				routePluginConfigs := RoutePluginConfigs{
					ID:   pluginInfo.ResID,
					Name: pluginInfo.Name,
					Tag:  pluginInfo.Key,
					Icon: pluginInfo.Icon,
					Type: pluginInfo.Type,
				}

				routePluginConfigInfo, ok := routesPluginConfigMap[routesPluginInfo.ID][pluginInfo.ResID]
				if ok {
					routePluginConfigs.Config = routePluginConfigInfo.Config
					routePluginConfigs.IsEnable = routePluginConfigInfo.IsEnable
					routePluginConfigs.ReleaseStatus = routePluginConfigInfo.ReleaseStatus
				}

				routePluginsList = append(routePluginsList, routePluginConfigs)
			}
		}
		routePluginConfigList.Plugins = routePluginsList

		list = append(list, routePluginConfigList)
	}

	return
}

func (r *Routes) RouteUpdateName(id string, name string) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if (len(id) == 0) || (len(name) == 0) {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Update("route_name", name).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *Routes) RouteSwitchEnable(id string, enable int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	routeInfo, routeInfoErr := r.RouteInfoById(id)
	if routeInfoErr != nil {
		return routeInfoErr
	}

	releaseStatus := routeInfo.ReleaseStatus
	if routeInfo.ReleaseStatus == utils.ReleaseStatusY {
		releaseStatus = utils.ReleaseStatusT
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(Routes{
			IsEnable:      enable,
			ReleaseStatus: releaseStatus}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *Routes) RouteSwitchRelease(id string, releaseStatus int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Update("release_status", releaseStatus).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}
