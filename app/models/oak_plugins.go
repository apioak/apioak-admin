package models

import (
	"apioak-admin/app/packages"
	"errors"
	"gorm.io/gorm"
)

type Plugins struct {
	ID          int    `gorm:"column:id;primary_key"` //primary key
	ResID       string `gorm:"column:res_id"`         //Plugin id
	PluginKey   string `gorm:"column:plugin_key"`     //Plugin key
	Icon        string `gorm:"column:icon"`           //Plugin icon
	Type        int    `gorm:"column:type"`           //Plugin type
	Description string `gorm:"column:description"`    //Plugin description
	ModelTime
}

// TableName sets the insert table name for this struct type
func (p *Plugins) TableName() string {
	return "oak_plugins"
}

func (p *Plugins) PluginAdd(pluginData *Plugins) error {
	err := packages.GetDb().
		Table(p.TableName()).
		Create(pluginData).Error

	return err
}

func (p *Plugins) PluginInfoByResId(resId string) (Plugins, error) {
	var pluginInfo Plugins
	err := packages.GetDb().Table(p.TableName()).
		Where("res_id = ?", resId).
		First(&pluginInfo).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return Plugins{}, err
	}

	return pluginInfo, nil
}

func (p *Plugins) PluginInfosByResIds(resIds []string) ([]Plugins, error) {
	pluginInfos := make([]Plugins, 0)
	err := packages.GetDb().
		Table(p.TableName()).
		Where("res_id IN ?", resIds).
		Find(&pluginInfos).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return pluginInfos, err
}

func (p *Plugins) PluginUpdate(resId string, pluginInfo *Plugins) error {
	updateError := packages.GetDb().
		Table(p.TableName()).
		Where("res_id = ?", resId).
		Updates(pluginInfo).Error

	return updateError
}

func (p *Plugins) PluginAllList() (list []Plugins, err error) {
	err = packages.GetDb().
		Table(p.TableName()).
		Find(&list).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (p *Plugins) PluginDelByPluginKeys(pluginKeys []string, filterResIds []string) (err error) {
	if len(pluginKeys) == 0 {
		return
	}

	dbObj := packages.GetDb().
		Table(p.TableName())
	if len(pluginKeys) != 0 {
		dbObj = dbObj.Where("plugin_key in ?", pluginKeys)
	}
	if len(filterResIds) != 0 {
		dbObj = dbObj.Where("res_id not in ?", filterResIds)
	}
	err = dbObj.Delete(p).Error

	return
}
