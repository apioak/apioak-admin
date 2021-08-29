package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
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
	IPTypeV4 = 1
	IPTypeV6 = 2
	sNodeId  = ""
)

// TableName sets the insert table name for this struct type
func (s *ServiceNodes) TableName() string {
	return "oak_service_nodes"
}

func IPTypeMap() map[string]int {
	var ipTypeMap map[string]int
	ipTypeMap = make(map[string]int)

	ipTypeMap[utils.IPV4] = IPTypeV4
	ipTypeMap[utils.IPV6] = IPTypeV6

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

	result := packages.GetDb().Table(s.TableName()).Select("id").First(&s)
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
	nodeInfos := []ServiceNodes{}
	packages.GetDb().Table(s.TableName()).Where("service_id IN ?", serviceIds).Find(&nodeInfos)

	return nodeInfos
}
