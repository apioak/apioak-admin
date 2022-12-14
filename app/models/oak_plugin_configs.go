package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"gorm.io/gorm"
)

const (
	PluginConfigsTypeService int = 1 // service
	PluginConfigsTypeRouter  int = 2 // router
)

type PluginConfigs struct {
	ID          int    `gorm:"column:id;primary_key"` // primary key
	ResID       string `gorm:"column:res_id"`         // Plugin config id
	Name        string `gorm:"column:name"`           // Plugin config name
	Type        int    `gorm:"column:type"`           // Plugin relation type 1:service  2:router
	TargetID    string `gorm:"column:target_id"`      // Target id
	PluginResID string `gorm:"column:plugin_res_id"`  // Plugin res id
	PluginKey   string `gorm:"column:plugin_key"`     // Plugin key
	Config      string `gorm:"column:config"`         // Plugin configuration
	Enable      int    `gorm:"column:enable"`         // Plugin config enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (m *PluginConfigs) TableName() string {
	return "oak_plugin_configs"
}

var recursionTimesPluginConfig = 1

func (m *PluginConfigs) ModelUniqueId() (string, error) {
	generateResId, generateIdErr := utils.IdGenerate(utils.IdTypeRouter)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateResId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesPluginConfig = 1
		return generateResId, nil
	} else {
		if recursionTimesPluginConfig == utils.IdGenerateMaxTimes {
			recursionTimesPluginConfig = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesPluginConfig++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (m *PluginConfigs) PluginConfigList(pluginConfigType int, targetId string) ([]PluginConfigs, error) {

	var pluginConfigs []PluginConfigs

	err := packages.GetDb().Table(m.TableName()).
		Where("type = ? AND target_id = ?", pluginConfigType, targetId).
		Find(&pluginConfigs).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return []PluginConfigs{}, err
	}

	return pluginConfigs, nil
}

func (m *PluginConfigs) PluginConfigInfoByResId(resId string) (PluginConfigs, error) {
	var pluginConfigInfo PluginConfigs
	err := packages.GetDb().Table(m.TableName()).
		Where("res_id = ?", resId).
		First(&pluginConfigInfo).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return PluginConfigs{}, err
	}

	return pluginConfigInfo, nil
}

func pluginConfigSyncTargetRelease(tx *gorm.DB, pluginConfigType int, pluginConfigTargetId string) error {
	if pluginConfigType == PluginConfigsTypeService {

		var service Services
		err := tx.Model(&Services{}).Where("res_id = ?", pluginConfigTargetId).First(&service).Error

		if err != nil {
			packages.Log.Error("Failed to modify the service plugin to obtain the service")
			return err
		}

		if service.Release == utils.ReleaseStatusY {
			service.Release = utils.ReleaseStatusT

			err = tx.Model(&Services{}).Updates(&service).Error

			if err != nil {
				packages.Log.Error("Failed to modify the service plugin release status")
				return err
			}
		}
	} else if pluginConfigType == PluginConfigsTypeRouter {
		var router Routers
		err := tx.Model(&Routers{}).Where("res_id = ?", pluginConfigTargetId).First(&router).Error

		if err != nil {
			packages.Log.Error("Failed to modify the router plugin to obtain the router")
			return err
		}

		if router.Release == utils.ReleaseStatusY {
			router.Release = utils.ReleaseStatusT

			err = tx.Model(&Services{}).Updates(&router).Error

			if err != nil {
				packages.Log.Error("Failed to modify the router plugin release status")
				return err
			}
		}
	}

	return nil
}

func (m *PluginConfigs) PluginConfigAdd(pluginConfigInfo *PluginConfigs) (string, error) {

	pluginConfigId, err := m.ModelUniqueId()

	if err != nil {
		return pluginConfigId, err
	}

	pluginConfigInfo.ResID = pluginConfigId
	if pluginConfigInfo.Name == "" {
		pluginConfigInfo.Name = pluginConfigId
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {

		err = packages.GetDb().Table(m.TableName()).Create(pluginConfigInfo).Error

		if err != nil {
			return err
		}

		err = pluginConfigSyncTargetRelease(tx, pluginConfigInfo.Type, pluginConfigInfo.TargetID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return pluginConfigId, err
	}

	return pluginConfigId, nil
}

func (m *PluginConfigs) PluginConfigUpdateColumns(
	pluginConfigId string,
	pluginConfigType int,
	pluginConfigTargetId string,
	params map[string]interface{},
) error {

	err := packages.GetDb().Transaction(func(tx *gorm.DB) error {

		err := tx.Table(m.TableName()).Where("res_id = ?", pluginConfigId).Updates(params).Error
		if err != nil {
			return err
		}

		err = pluginConfigSyncTargetRelease(tx, pluginConfigType, pluginConfigTargetId)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *PluginConfigs) PluginConfigDelete(
	pluginConfigId string,
	pluginConfigType int,
	pluginConfigTargetId string,
) error {

	err := packages.GetDb().Transaction(func(tx *gorm.DB) error {

		err := tx.Table(m.TableName()).Where("res_id = ?", pluginConfigId).Delete(&PluginConfigs{}).Error
		if err != nil {
			return err
		}

		err = pluginConfigSyncTargetRelease(tx, pluginConfigType, pluginConfigTargetId)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
