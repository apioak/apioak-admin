package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
)

type ServiceNodes struct {
	ID         string `gorm:"column:id;primary_key"` //Service node id
	ServiceID  string `gorm:"column:service_id"`     //Service id
	NodeIP     string `gorm:"column:node_ip"`        //Node IP
	IPType     int    `gorm:"column:ip_type"`        //IP Type  1:IPV4  2:IPV6
	NodePort   int    `gorm:"column:node_port"`      //Node port
	NodeWeight int    `gorm:"column:node_weight"`    //Node weight
	ModelTime
}

var (
	sNodeId = ""
)

// TableName sets the insert table name for this struct type
func (s *ServiceNodes) TableName() string {
	return "oak_service_nodes"
}

var recursionTimesServiceNodes = 1

func (m *ServiceNodes) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeServiceNode)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("id = ?", generateId).
		Select("id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesServiceNodes = 1
		return generateId, nil
	} else {
		if recursionTimesServiceNodes == utils.IdGenerateMaxTimes {
			recursionTimesServiceNodes = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesServiceNodes++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func IPTypeMap() map[string]int {
	var ipTypeMap map[string]int
	ipTypeMap = make(map[string]int)

	ipTypeMap[utils.IPV4] = utils.IPTypeV4
	ipTypeMap[utils.IPV6] = utils.IPTypeV6

	return ipTypeMap
}

func (s *ServiceNodes) ServiceNodeIdUnique(sNodeIds map[string]string) (string, error) {
	if s.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeServiceNode)
		if err != nil {
			return "", err
		}
		s.ID = tmpID
	}

	result := packages.GetDb().
		Table(s.TableName()).
		Select("id").
		First(&s)

	mapId := sNodeIds[s.ID]
	if (result.RowsAffected == 0) && (s.ID != mapId) {
		sNodeId = s.ID
		sNodeIds[s.ID] = s.ID
		return sNodeId, nil
	} else {
		svcNodeId, svcErr := utils.IdGenerate(utils.IdTypeServiceNode)
		if svcErr != nil {
			return "", svcErr
		}
		s.ID = svcNodeId
		_, err := s.ServiceNodeIdUnique(sNodeIds)
		if err != nil {
			return "", err
		}
	}

	return sNodeId, nil
}

func (s *ServiceNodes) NodeInfosByServiceIds(serviceIds []string) []ServiceNodes {
	nodeInfos := make([]ServiceNodes, 0)
	packages.GetDb().
		Table(s.TableName()).
		Where("service_id IN ?", serviceIds).
		Find(&nodeInfos)

	return nodeInfos
}
