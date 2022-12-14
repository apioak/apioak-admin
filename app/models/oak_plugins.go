package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
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

var recursionTimesPlugins = 1

func (m *Plugins) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypePlugin)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesPlugins = 1
		return generateId, nil
	} else {
		if recursionTimesPlugins == utils.IdGenerateMaxTimes {
			recursionTimesPlugins = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesPlugins++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
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

func (p *Plugins) PluginAllList() []Plugins {
	pluginAllList := make([]Plugins, 0)
	packages.GetDb().
		Table(p.TableName()).
		Find(&pluginAllList)

	return pluginAllList
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
