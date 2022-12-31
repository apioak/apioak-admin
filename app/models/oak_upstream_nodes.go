package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"gorm.io/gorm"
)

type UpstreamNodes struct {
	ID            int    `gorm:"column:id;primary_key"`  // primary key
	ResID         string `gorm:"column:res_id"`          // Service node id
	UpstreamResID string `gorm:"column:upstream_res_id"` // Upstream id
	NodeIP        string `gorm:"column:node_ip"`         // Node IP
	IPType        int    `gorm:"column:ip_type"`         // IP Type  1:IPV4  2:IPV6
	NodePort      int    `gorm:"column:node_port"`       // Node port
	NodeWeight    int    `gorm:"column:node_weight"`     // Node weight
	Health        int    `gorm:"column:health"`          // Health type  1:HEALTH  2:UNHEALTH
	HealthCheck   int    `gorm:"column:health_check"`    // Health check  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *UpstreamNodes) TableName() string {
	return "oak_upstream_nodes"
}

func (m *UpstreamNodes) ModelUniqueId() (generateId string, err error) {
	generateId, err = utils.IdGenerate(utils.IdTypeUpstreamNode)
	if err != nil {
		return
	}

	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err == nil {
		recursionTimesServices = 1
		return
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			err = errors.New(enums.CodeMessages(enums.IdConflict))
			return
		}

		recursionTimesServices++
		generateId, err = m.ModelUniqueId()
		if err != nil {
			return
		}

		return
	}
}

func (m *UpstreamNodes) UpstreamNodeListByResIds(upstreamNodeResIds []string) (list []UpstreamNodes, err error) {
	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id in ?", upstreamNodeResIds).
		Find(&list).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (m *UpstreamNodes) UpstreamNodeListByUpstreamResIds(upstreamResIds []string) (list []UpstreamNodes, err error) {
	err = packages.GetDb().
		Table(m.TableName()).
		Where("upstream_res_id in ?", upstreamResIds).
		Find(&list).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}
