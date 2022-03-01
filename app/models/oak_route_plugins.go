package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
)

type RoutePlugins struct {
	ID            string `gorm:"column:id;primary_key"`      //Plugin id
	RouteID       string `gorm:"column:route_id;primaryKey"` //Route id
	PluginID      string `gorm:"column:plugin_id"`           //Plugin id
	Order         int    `gorm:"column:order"`               //Order sort
	Config        string `gorm:"column:config"`              //Routing configuration
	IsEnable      int    `gorm:"column:is_enable"`           //Plugin enable  1:on  2:off
	ReleaseStatus int    `gorm:"column:release_status"`      //Route plugin release status 1:unpublished  2:to be published  3:published
	ModelTime
	Plugin Plugins `gorm:"foreignKey:PluginID;"`
}

// TableName sets the insert table name for this struct type
func (r *RoutePlugins) TableName() string {
	return "oak_route_plugins"
}

var recursionTimesRoutePlugins = 1

func (m *RoutePlugins) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeRoutePlugin)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("id = ?", generateId).
		Select("id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesRoutePlugins = 1
		return generateId, nil
	} else {
		if recursionTimesRoutePlugins == utils.IdGenerateMaxTimes {
			recursionTimesRoutePlugins = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesRoutePlugins++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (r *RoutePlugins) RoutePluginInfosByPluginIds(pluginIds []string) ([]RoutePlugins, error) {
	routePluginInfos := make([]RoutePlugins, 0)
	if len(pluginIds) == 0 {
		return routePluginInfos, nil
	}

	err := packages.GetDb().
		Table(r.TableName()).
		Where("plugin_id IN ?", pluginIds).
		Order("created_at DESC").
		Find(&routePluginInfos).Error

	return routePluginInfos, err
}

func (r *RoutePlugins) RoutePluginAllListByRouteIds(routeIds []string) []RoutePlugins {
	routePluginAllList := make([]RoutePlugins, 0)
	packages.GetDb().
		Table(r.TableName()).
		Where("route_id IN ?", routeIds).
		Order("created_at DESC").
		Find(&routePluginAllList)

	return routePluginAllList
}

func (r *RoutePlugins) RoutePluginInfoConfigListByRouteIds(routeIds []string) []RoutePlugins {
	routePluginInfoConfigList := make([]RoutePlugins, 0)
	packages.GetDb().
		Table(r.TableName()).
		Where("route_id IN ?", routeIds).
		Order("created_at DESC").
		Preload("Plugin").
		Find(&routePluginInfoConfigList)

	return routePluginInfoConfigList
}

func (r *RoutePlugins) RoutePluginInfoById(id string, routeId string, pluginId string) RoutePlugins {
	routePluginInfo := RoutePlugins{}
	db := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id)

	if len(routeId) != 0 {
		db = db.Where("route_id = ?", routeId)
	}
	if len(pluginId) != 0 {
		db = db.Where("plugin_id = ?", pluginId)
	}

	db.First(&routePluginInfo)

	return routePluginInfo
}

func (r *RoutePlugins) RoutePluginInfoByRoutePluginId(routeId string, pluginId string) RoutePlugins {
	routePluginInfo := RoutePlugins{}
	packages.GetDb().
		Table(r.TableName()).
		Where("route_id = ?", routeId).
		Where("plugin_id = ?", pluginId).
		First(&routePluginInfo)

	return routePluginInfo
}

func (r *RoutePlugins) RoutePluginInfosByRouteIdRelease(routeIds []string, releaseStatus []int) []RoutePlugins {
	routePluginInfos := make([]RoutePlugins, 0)
	if len(routeIds) == 0 {
		return routePluginInfos
	}

	db := packages.GetDb().
		Table(r.TableName()).
		Where("route_id IN ?", routeIds)
	if len(releaseStatus) != 0 {
		db = db.Where("release_status IN ?", releaseStatus)
	}
	db.Find(&routePluginInfos)

	return routePluginInfos
}

func (r *RoutePlugins) RoutePluginAdd(routePluginData *RoutePlugins) error {
	routePluginId, routePluginIdUniqueErr := r.ModelUniqueId()
	if routePluginIdUniqueErr != nil {
		return routePluginIdUniqueErr
	}
	routePluginData.ID = routePluginId

	err := packages.GetDb().
		Table(r.TableName()).
		Create(routePluginData).Error

	return err
}

func (r *RoutePlugins) RoutePluginUpdateColumnsByIds(ids []string, routePluginData *RoutePlugins) error {
	return packages.GetDb().
		Table(r.TableName()).
		Where("id IN ?", ids).
		Updates(routePluginData).Error
}

func (r *RoutePlugins) RoutePluginUpdate(id string, routePluginData *RoutePlugins) error {
	err := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(routePluginData).Error

	return err
}

func (r *RoutePlugins) RoutePluginSwitchEnable(id string, enable int) error {
	routePluginInfo := r.RoutePluginInfoById(id, "", "")
	releaseStatus := routePluginInfo.ReleaseStatus
	if routePluginInfo.ReleaseStatus == utils.ReleaseStatusY {
		releaseStatus = utils.ReleaseStatusT
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(RoutePlugins{
			IsEnable:      enable,
			ReleaseStatus: releaseStatus}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *RoutePlugins) RoutePluginSwitchRelease(id string, release int) error {
	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Update("release_status", release).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *RoutePlugins) RoutePluginDelete(id string) error {
	deleteErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Delete(r).Error

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (r *RoutePlugins) RoutePluginInfosByRouteId(routeId string) []RoutePlugins {
	routePluginInfos := make([]RoutePlugins, 0)
	packages.GetDb().
		Table(r.TableName()).
		Where("route_id = ?", routeId).
		Find(&routePluginInfos)

	return routePluginInfos
}
