package models

import "apioak-admin/app/packages"

type RoutePlugins struct {
	ID       string `gorm:"column:id;primary_key"`      //Plugin id
	RouteID  string `gorm:"column:route_id;primaryKey"` //Route id
	PluginID string `gorm:"column:plugin_id"`           //Plugin id
	Order    int    `gorm:"column:order"`               //Order sort
	Config   string `gorm:"column:config"`              //Routing configuration
	IsEnable int    `gorm:"column:is_enable"`           //Plugin enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (r *RoutePlugins) TableName() string {
	return "oak_route_plugins"
}

func (r *RoutePlugins) RoutePluginInfosByPluginIds(pluginIds []string) ([]RoutePlugins, error) {
	routePluginInfos := make([]RoutePlugins, 0)
	if len(pluginIds) == 0 {
		return routePluginInfos, nil
	}

	err := packages.GetDb().Table(r.TableName()).Where("plugin_id IN ?", pluginIds).Find(&routePluginInfos).Error

	return routePluginInfos, err
}
