package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
)

type RoutePlugins struct {
	ID        string `gorm:"column:id;primary_key"`      //Plugin id
	RouteID   string `gorm:"column:route_id;primaryKey"` //Route id
	PluginID  string `gorm:"column:plugin_id"`           //Plugin id
	Order     int    `gorm:"column:order"`               //Order sort
	Config    string `gorm:"column:config"`              //Routing configuration
	IsEnable  int    `gorm:"column:is_enable"`           //Plugin enable  1:on  2:off
	IsRelease int    `gorm:"column:is_release"`          //Route plugin release  1:on  2:off
	ModelTime
	Plugin Plugins `gorm:"foreignKey:PluginID;"`
}

// TableName sets the insert table name for this struct type
func (r *RoutePlugins) TableName() string {
	return "oak_route_plugins"
}

var routePluginId = ""

func (r *RoutePlugins) PluginIdUnique(rIds map[string]string) (string, error) {
	if r.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeRoutePlugin)
		if err != nil {
			return "", err
		}
		r.ID = tmpID
	}

	result := packages.GetDb().
		Table(r.TableName()).
		Select("id").
		First(&r)

	mapId := rIds[r.ID]
	if (result.RowsAffected == 0) && (r.ID != mapId) {
		routePluginId = r.ID
		rIds[r.ID] = r.ID
		return routePluginId, nil
	} else {
		rpId, rpIdErr := utils.IdGenerate(utils.IdTypeRoutePlugin)
		if rpIdErr != nil {
			return "", rpIdErr
		}
		r.ID = rpId
		_, err := r.PluginIdUnique(rIds)
		if err != nil {
			return "", err
		}
	}

	return routePluginId, nil
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

func (r *RoutePlugins) RoutePluginAdd(routePluginData *RoutePlugins) error {
	tpmIds := map[string]string{}
	routePluginId, routePluginIdUniqueErr := r.PluginIdUnique(tpmIds)
	if routePluginIdUniqueErr != nil {
		return routePluginIdUniqueErr
	}
	routePluginData.ID = routePluginId

	err := packages.GetDb().
		Table(r.TableName()).
		Create(routePluginData).Error

	return err
}

func (r *RoutePlugins) RoutePluginUpdate(id string, routePluginData *RoutePlugins) error {
	err := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(routePluginData).Error

	return err
}

func (r *RoutePlugins) RoutePluginSwitchEnable(id string, enable int) error {
	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Updates(RoutePlugins{IsEnable: enable, IsRelease: utils.IsReleaseN}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *RoutePlugins) RoutePluginSwitchRelease(id string, release int) error {
	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("id = ?", id).
		Update("is_release", release).Error

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
