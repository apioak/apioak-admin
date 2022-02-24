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

func (r *Routes) RouteInfosById(routeId string) (Routes, error) {
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

func (r *Routes) RouteCopy(routeData Routes, routePlugins []RoutePlugins) (string, []string, error) {
	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	routePluginIds := make([]string, 0)
	if err := tx.Error; err != nil {
		return "", routePluginIds, err
	}

	routeId, routeIdUniqueErr := r.ModelUniqueId()
	if routeIdUniqueErr != nil {
		return routeId, routePluginIds, routeIdUniqueErr
	}

	routeData.ID = routeId
	routeData.RouteName = routeId

	addErr := tx.
		Table(r.TableName()).
		Create(&routeData).Error

	if addErr != nil {
		tx.Rollback()
		return routeId, routePluginIds, addErr
	}

	if len(routePlugins) != 0 {
		routePluginModel := RoutePlugins{}
		for k, _ := range routePlugins {
			routePluginId, routePluginIdErr := routePluginModel.ModelUniqueId()
			if routePluginIdErr != nil {
				tx.Rollback()
				return routeId, routePluginIds, routePluginIdErr
			}

			routePlugins[k].ID = routePluginId
			routePlugins[k].RouteID = routeId
			routePlugins[k].CreatedAt = time.Now()
			routePlugins[k].UpdatedAt = time.Now()

			routePluginIds = append(routePluginIds, routePluginId)
		}

		addRoutePluginErr := tx.
			Table(routePluginModel.TableName()).
			Create(&routePlugins).Error

		if addRoutePluginErr != nil {
			tx.Rollback()
			return routeId, routePluginIds, addRoutePluginErr
		}
	}

	return routeId, routePluginIds, tx.Commit().Error
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

func (r *Routes) RouteListPage(
	serviceId string,
	param *validators.ValidatorRouteList) (list []Routes, total int, listError error) {
	tx := packages.GetDb().
		Table(r.TableName()).
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

	listError = ListPaginate(tx, &list, &param.BaseListPage)

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

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(Routes{
			IsEnable:      enable,
			ReleaseStatus: utils.ReleaseStatusT}).Error

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