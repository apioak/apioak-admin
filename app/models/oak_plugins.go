package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type Plugins struct {
	ID          int       `gorm:"column:id;primary_key"` //primary key
	ResID       string    `gorm:"column:res_id"`         //Plugin id
	Name        string    `gorm:"column:name"`           //Plugin name
	Key         string    `gorm:"column:key"`            //Plugin key
	Icon        string    `gorm:"column:icon"`           //Plugin icon
	Type        int       `gorm:"column:type"`           //Plugin type
	Description string    `gorm:"column:description"`    //Plugin description
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

func (p *Plugins) PluginInfosByKeys(keys []string, filterPluginIds []string) ([]Plugins, error) {
	pluginInfos := make([]Plugins, 0)
	db := packages.GetDb().
		Table(p.TableName()).
		Where("key IN ?", keys)

	if len(filterPluginIds) != 0 {
		db = db.Where("res_id NOT IN ?", filterPluginIds)
	}

	err := db.Find(&pluginInfos).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return pluginInfos, err
}

func (p *Plugins) PluginAdd(pluginData *Plugins) error {
	pluginResId, pluginIdUniqueErr := p.ModelUniqueId()
	if pluginIdUniqueErr != nil {
		return pluginIdUniqueErr
	}
	pluginData.ResID = pluginResId

	err := packages.GetDb().
		Table(p.TableName()).
		Create(pluginData).Error

	return err
}

func (p *Plugins) PluginInfoByResId(resId string) Plugins {
	pluginInfo := Plugins{}
	packages.GetDb().
		Table(p.TableName()).
		Where("res_id = ?", resId).
		Find(&pluginInfo)

	return pluginInfo
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

func (p *Plugins) PluginInfoByResIdRouteServiceId(pluginResId string) Plugins {
	pluginInfo := Plugins{}
	packages.GetDb().
		Table(p.TableName()).
		Where("res_id = ?", pluginResId).
		First(&pluginInfo)

	return pluginInfo
}

func (p *Plugins) PluginUpdate(id string, pluginInfo *Plugins) error {
	updateError := packages.GetDb().
		Table(p.TableName()).
		Where("res_id = ?", id).
		Updates(pluginInfo).Error

	return updateError
}

func (p *Plugins) PluginDelete(resId string) error {
	deleteError := packages.GetDb().
		Table(p.TableName()).
		Where("res_id = ?", resId).
		Delete(p).Error

	return deleteError
}

func (p *Plugins) PluginListPage(param *validators.PluginList) (list []Plugins, total int, listError error) {
	tx := packages.GetDb().
		Table(p.TableName())

	if param.Type != 0 {
		tx = tx.Where("type = ?", param.Type)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		orWhere := packages.GetDb().
			Where("name LIKE ?", search).
			Or("Tag LIKE ?", search).
			Or("description LIKE ?", search)
		tx = tx.Where(orWhere)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("created_at desc")
	listError = ListPaginate(tx, &list, &param.BaseListPage)
	return
}

func (p *Plugins) PluginAllList() []Plugins {
	pluginAllList := make([]Plugins, 0)
	packages.GetDb().
		Table(p.TableName()).
		Find(&pluginAllList)

	return pluginAllList
}
