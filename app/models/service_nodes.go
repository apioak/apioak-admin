package models

import "apioak-admin/app/utils"

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
)

// TableName sets the insert table name for this struct type
func (s *ServiceNodes) TableName() string {
	return "oak_service_nodes"
}

func IPTypeMap() map[string]int {
	var ipTypeMap map[string]int
	ipTypeMap =  make(map[string]int)

	ipTypeMap[utils.IPV4] = IPTypeV4
	ipTypeMap[utils.IPV6] = IPTypeV6

	return ipTypeMap
}

