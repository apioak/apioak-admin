package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"strings"
)

type Plugins struct {
	ID          string `gorm:"column:id;primary_key"` //Plugin id
	Name        string `gorm:"column:name"`           //Plugin name
	Tag         string `gorm:"column:tag"`            //Plugin tag
	Icon        string `gorm:"column:icon"`           //Plugin icon
	Type        int    `gorm:"column:type"`           //Plugin type
	Description string `gorm:"column:description"`    //Plugin description
	ModelTime
}

// TableName sets the insert table name for this struct type
func (p *Plugins) TableName() string {
	return "oak_plugins"
}

var pluginId = ""

func (p *Plugins) PluginIdUnique(sIds map[string]string) (string, error) {
	if p.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypePlugin)
		if err != nil {
			return "", err
		}
		p.ID = tmpID
	}

	result := packages.GetDb().
		Table(p.TableName()).
		Select("id").
		First(&p)

	mapId := sIds[p.ID]
	if (result.RowsAffected == 0) && (p.ID != mapId) {
		pluginId = p.ID
		sIds[p.ID] = p.ID
		return pluginId, nil
	} else {
		svcId, svcErr := utils.IdGenerate(utils.IdTypePlugin)
		if svcErr != nil {
			return "", svcErr
		}
		p.ID = svcId
		_, err := p.PluginIdUnique(sIds)
		if err != nil {
			return "", err
		}
	}

	return pluginId, nil
}

func (p *Plugins) PluginInfosByTags(tag []string, filterPluginIds []string) ([]Plugins, error) {
	pluginInfos := make([]Plugins, 0)
	db := packages.GetDb().
		Table(p.TableName()).
		Where("tag IN ?", tag)

	if len(filterPluginIds) != 0 {
		db = db.Where("id NOT IN ?", filterPluginIds)
	}

	err := db.Find(&pluginInfos).Error

	return pluginInfos, err
}

func (p *Plugins) PluginAdd(pluginData *Plugins) error {
	tpmIds := map[string]string{}
	pluginId, pluginIdUniqueErr := p.PluginIdUnique(tpmIds)
	if pluginIdUniqueErr != nil {
		return pluginIdUniqueErr
	}
	pluginData.ID = pluginId

	err := packages.GetDb().
		Table(p.TableName()).
		Create(pluginData).Error

	return err
}

func (p *Plugins) PluginInfosByIds(ids []string) ([]Plugins, error) {
	pluginInfos := make([]Plugins, 0)
	err := packages.GetDb().
		Table(p.TableName()).
		Where("id IN ?", ids).
		Find(&pluginInfos).Error

	return pluginInfos, err
}

func (p *Plugins) PluginInfoByIdRouteServiceId(pluginId string) Plugins {
	pluginInfo := Plugins{}
	packages.GetDb().
		Table(p.TableName()).
		Where("id = ?", pluginId).
		First(&pluginInfo)

	return pluginInfo
}

func (p *Plugins) PluginUpdate(id string, pluginInfo *Plugins) error {
	updateError := packages.GetDb().
		Table(p.TableName()).
		Where("id = ?", id).
		Updates(pluginInfo).Error

	return updateError
}

func (p *Plugins) PluginDelete(id string) error {
	deleteError := packages.GetDb().
		Table(p.TableName()).
		Where("id = ?", id).
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
