package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
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

func (m *UpstreamNodes) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeUpstreamNode)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesServices = 1
		return generateId, nil
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesServices++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}
